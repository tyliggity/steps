package main

import (
	"fmt"

	"github.com/google/shlex"
	"github.com/stackpulse/steps-sdk-go/step"
	"github.com/stackpulse/steps/kubectl/base"
)

type podFiles struct {
	localPath  string
	remotePath string
	pod        string
	container  string
	namespace  string
}

func (p podFiles) command(isDownload bool) string {
	result := "cp "

	podPath := fmt.Sprintf("%s:%s", p.pod, p.remotePath)
	if p.namespace != "" {
		podPath = fmt.Sprintf("%s/%s", p.namespace, podPath)
	}

	if isDownload {
		result += fmt.Sprintf("%s %s", podPath, p.localPath)
	} else {
		result += fmt.Sprintf("%s %s", p.localPath, podPath)
	}

	if p.container != "" {
		result += " -c " + p.container
	}

	return result
}

func (p podFiles) copy(kbl *base.KubectlStep, isDownload bool) ([]byte, int, error) {
	copyCmd := p.command(isDownload)
	cmdArgs, err := shlex.Split(copyCmd)
	if err != nil {
		return nil, step.ExitCodeFailure, fmt.Errorf("split copy command: %w", err)
	}

	kbl.StepArgs.Format = "default"
	return kbl.Execute(cmdArgs)
}

func (p podFiles) upload(kbl *base.KubectlStep) ([]byte, int, error) {
	return p.copy(kbl, false)
}

func (p podFiles) download(kbl *base.KubectlStep) ([]byte, int, error) {
	return p.copy(kbl, true)
}
