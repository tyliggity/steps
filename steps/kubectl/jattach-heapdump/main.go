package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"
	"strings"
	"time"

	"github.com/stackpulse/steps-sdk-go/env"
	"github.com/stackpulse/steps-sdk-go/step"
	"github.com/stackpulse/steps/kubectl/base"
	_ "gocloud.dev/blob/gcsblob"
	_ "gocloud.dev/blob/s3blob"
)

type memoryDump struct {
	base.Args
	Storage        storage
	JAttach        jattach
	PodName        string `env:"POD_NAME,required"`
	Container      string `env:"CONTAINER"`
	Namespace      string `env:"NAMESPACE"`
	ContainerShell string `env:"CONTAINER_SHELL" envDefault:"bash"`
	kubectl        *base.KubectlStep
}

type output struct {
	DumpFileURI string `json:"dump_file_uri"`
	step.Outputs
}

func closeWithLog(closer io.Closer) {
	if err := closer.Close(); err != nil {
		log.Println("Failed to perform close: ", err)
	}
}

func (m *memoryDump) Init() error {
	err := env.Parse(m)
	if err != nil {
		return err
	}

	kubectl, err := base.NewKubectlStep(&m.Args, true)
	if err != nil {
		return fmt.Errorf("parse kubectl step params: %w", err)
	}

	err = m.Storage.initGCloud()
	if err != nil {
		return fmt.Errorf("init GCloud credentials: %w", err)
	}

	m.kubectl = kubectl

	return nil
}

func (m *memoryDump) Run() (int, []byte, error) {
	var err error
	defer func() {
		if cleanErr := m.cleanup(); cleanErr != nil {
			log.Println("Cleanup failed due to: ", cleanErr)
		}
	}()
	err = m.prerequisites()
	if err != nil {
		return step.ExitCodeFailure, nil, fmt.Errorf("install prerequisites: %w", err)
	}

	err = m.dump()
	if err != nil {
		return step.ExitCodeFailure, nil, fmt.Errorf("dump heap memory: %w", err)
	}

	err = m.pull()
	if err != nil {
		return step.ExitCodeFailure, nil, fmt.Errorf("pull dump file from pod: %w", err)
	}

	dumpFileKey := fmt.Sprintf("heapdump_%d.hprof", time.Now().Unix())
	err = m.Storage.upload(context.Background(), map[string]string{
		pulledDumpLocalPath: dumpFileKey,
	})
	if err != nil {
		return step.ExitCodeFailure, nil, fmt.Errorf("upload dump file: %w", err)
	}

	jsonOutput, err := json.Marshal(&output{
		DumpFileURI: m.Storage.Bucket + "/" + dumpFileKey,
	})
	if err != nil {
		return step.ExitCodeFailure, nil, fmt.Errorf("marshal output: %w", err)
	}

	env.SetFormatter(env.JsonFormat, true)
	return step.ExitCodeOK, jsonOutput, nil
}

func (m *memoryDump) prerequisites() error {
	dumpScriptFile, err := os.Create(dumpScriptLocalPath)
	if err != nil {
		return fmt.Errorf("open dump script file at '%s': %w", dumpScriptLocalPath, err)
	}
	defer closeWithLog(dumpScriptFile)

	err = m.JAttach.WriteDumpScript(dumpScriptFile)
	if err != nil {
		return fmt.Errorf("generate dump script: %w", err)
	}

	// upload required files
	for local, remote := range map[string]string{
		localJAttachPath:    m.JAttach.RemotePath,
		dumpScriptLocalPath: m.JAttach.DumpScriptRemotePath,
	} {
		out, _, err := podFiles{
			localPath:  local,
			remotePath: remote,
			pod:        m.PodName,
			container:  m.Container,
			namespace:  m.Namespace,
		}.upload(m.kubectl)
		if err != nil {
			return fmt.Errorf("upload file '%s' to '%s' due to %s: %w", local, remote, string(out), err)
		}
		m.kubectl.Debugln(string(out))
	}

	return nil
}

func (m *memoryDump) exec(command string) error {
	baseCmd, err := m.kubectl.BaseCommand()
	if err != nil {
		return fmt.Errorf("kubectl base command: %w", err)
	}
	execCmd := append(append([]string{"exec", m.PodName}, baseCmd...), "--", m.ContainerShell, "-c", command)

	out, _, err := m.kubectl.Execute(execCmd)
	if err != nil {
		return fmt.Errorf("execute command due to %s: %w", string(out), err)
	}
	m.kubectl.Debugln(string(out))
	return nil
}

func (m *memoryDump) dump() error {
	return m.exec(
		fmt.Sprintf("chmod +x %s && %s", m.JAttach.DumpScriptRemotePath, m.JAttach.DumpScriptRemotePath),
	)
}

func (m *memoryDump) pull() error {
	for local, remote := range map[string]string{
		m.JAttach.OutputDumpPath: pulledDumpLocalPath,
	} {
		out, _, err := podFiles{
			localPath:  local,
			remotePath: remote,
			pod:        m.PodName,
			container:  m.Container,
			namespace:  m.Namespace,
		}.download(m.kubectl)
		if err != nil {
			return fmt.Errorf("download file from '%s' to '%s' due to %s: %w", remote, local, string(out), err)
		}
		m.kubectl.Debugln(string(out))
	}

	return nil
}

func (m *memoryDump) cleanup() error {
	files := []string{
		m.JAttach.DumpScriptRemotePath,
		m.JAttach.OutputDumpPath,
		m.JAttach.RemotePath,
	}
	return m.exec(fmt.Sprintf("rm -f %s", strings.Join(files, " ")))
}

func main() {
	step.Run(&memoryDump{})
}
