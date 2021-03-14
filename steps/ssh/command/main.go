package main

import (
	"fmt"
	envconf "github.com/caarlos0/env/v6"
	"github.com/stackpulse/public-steps/common/env"
	"github.com/stackpulse/public-steps/common/exec"
	"github.com/stackpulse/public-steps/common/step"
	"io/ioutil"
	"strconv"
	"time"
)

type Args struct {
	Username              string        `env:"USERNAME,required" envDefault:""`
	Hostname              string        `env:"HOSTNAME,required" envDefault:""`
	Command               string        `env:"COMMAND,required" envDefault:""`
	PrivateKey            string        `env:"PRIVATE_KEY,required" envDefault:""`
	StrictHostKeyChecking string        `env:"STRICT_HOST_KEY_CHECKING" envDefault:"no"`
	LogLevel              string        `env:"LOG_LEVEL" envDefault:"ERROR"`
	Port                  int           `env:"PORT" envDefault:"22"`
	ConnectionTimeout     time.Duration `env:"CONNECTION_TIMEOUT" envDefault:"30s"`
}

const (
	PrivateKeyPath = "/key.pem"
)

type SSHCommand struct {
	args Args
}

func (s *SSHCommand) Init() error {
	err := envconf.Parse(&s.args)
	if err != nil {
		return err
	}

	return nil
}

func (s *SSHCommand) Run() (int, []byte, error) {
	env.SetFormatter(env.RawFormat, true)
	err := ioutil.WriteFile(PrivateKeyPath, []byte(s.args.PrivateKey), 0644)
	if err != nil {
		return step.ExitCodeFailure, nil, fmt.Errorf("run ssh command: write file: error: %w", err)
	}

	// Restrict key.pem file capabilities (for ssh usage)
	output, exitCode, err := exec.Execute("chmod", []string{"600", PrivateKeyPath})
	if err != nil {
		return exitCode, nil, fmt.Errorf("execute linux command chmod: err code: %d, error: %w", exitCode, err)
	}

	sshArgs := s.buildArgs(s.args.Username, s.args.Hostname, s.args.Command, s.args.StrictHostKeyChecking, s.args.LogLevel, s.args.Port, s.args.ConnectionTimeout)
	output, exitCode, err = exec.Execute("ssh", sshArgs)
	if err != nil {
		return exitCode, nil, fmt.Errorf("execute linux command ssh: err code: %d, error: %w", exitCode, err)
	}

	return step.ExitCodeOK, output, err
}

func (s *SSHCommand) buildArgs(username, hostname, linuxCmd, StrictHostKeyChecking, LogLevel string, port int, connectionTimeout time.Duration) []string {
	args := []string{"-o", fmt.Sprintf("StrictHostKeyChecking=%s", StrictHostKeyChecking)}
	args = append(args, "-o", fmt.Sprintf("LogLevel=%s", LogLevel))
	args = append(args, "-i", PrivateKeyPath)
	args = append(args, fmt.Sprintf("%s@%s", username, hostname))
	args = append(args, "-p", strconv.Itoa(port))
	args = append(args, fmt.Sprintf("-oConnectTimeout=%d", int(connectionTimeout.Seconds())))
	args = append(args, linuxCmd)

	return args
}

func main() {
	step.Run(&SSHCommand{})
}
