package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"strconv"
	"time"

	envconf "github.com/caarlos0/env/v6"
	sdkExec "github.com/stackpulse/steps-sdk-go/exec"
	"github.com/stackpulse/steps-sdk-go/step"
)

type Args struct {
	Username              string        `env:"USERNAME,required" envDefault:""`
	Hostname              string        `env:"HOSTNAME,required" envDefault:""`
	Command               string        `env:"COMMAND,required" envDefault:""`
	PrivateKey            string        `env:"PRIVATE_KEY" envDefault:""`
	Password              string        `env:"PASSWORD" envDefault:""`
	StrictHostKeyChecking string        `env:"STRICT_HOST_KEY_CHECKING" envDefault:"no"`
	LogLevel              string        `env:"LOG_LEVEL" envDefault:"ERROR"`
	Port                  int           `env:"PORT" envDefault:"22"`
	ConnectionTimeout     time.Duration `env:"CONNECTION_TIMEOUT" envDefault:"30s"`
}

type Outputs struct {
	CommandOutput string `json:"output"`
	step.Outputs
}

const (
	PrivateKeyPath = "/key"
	LogFilePath    = "/log"
)

type SSHCommand struct {
	args Args
}

func (s *SSHCommand) Init() error {
	err := envconf.Parse(&s.args)
	if err != nil {
		return err
	}

	if s.args.PrivateKey == "" && s.args.Password == "" {
		return fmt.Errorf("private key or password is required")
	}

	return nil
}

// try to print ssh log, ignore errors.
func (s *SSHCommand) PrintLog() {
	content, _ := ioutil.ReadFile(LogFilePath)
	fmt.Println("--SSHLOG--")
	fmt.Printf(string(content))
	fmt.Println("----------")
}

func (s *SSHCommand) marshalOutput(output string) []byte {
	outputJSON, _ := json.Marshal(Outputs{
		CommandOutput: output,
	})
	return outputJSON
}

func (s *SSHCommand) Run() (int, []byte, error) {
	var sshCmd string
	var sshArgs []string
	if s.args.PrivateKey != "" {
		err := ioutil.WriteFile(PrivateKeyPath, []byte(s.args.PrivateKey), 0644)
		if err != nil {
			return step.ExitCodeFailure, s.marshalOutput(""), fmt.Errorf("write private key: %w", err)
		}

		// Restrict key.pem file capabilities (for ssh usage)
		output, exitCode, err := sdkExec.Execute("chmod", []string{"600", PrivateKeyPath})
		if err != nil {
			return exitCode, s.marshalOutput(""), fmt.Errorf("chmod private key: %w: %s", err, output)
		}

		sshCmd, sshArgs = s.buildPrivateKeyCommand(s.args.Username, s.args.Hostname, s.args.Command, s.args.StrictHostKeyChecking, s.args.LogLevel, s.args.Port, s.args.ConnectionTimeout)

	} else {
		sshCmd, sshArgs = s.buildPasswordCommand(s.args.Username, s.args.Hostname, s.args.Command, s.args.StrictHostKeyChecking, s.args.LogLevel, s.args.Port, s.args.ConnectionTimeout, s.args.Password)
	}

	output, exitCode, err := sdkExec.Execute(sshCmd, sshArgs)
	s.PrintLog()
	fmt.Println("--OUTPUT--")
	fmt.Printf(string(output))
	fmt.Println("----------")

	marshaledOuptut := s.marshalOutput(string(output))

	if err != nil {
		return exitCode, marshaledOuptut, fmt.Errorf("execute ssh: %w", err)
	}

	return step.ExitCodeOK, marshaledOuptut, nil
}

func (s *SSHCommand) buildPrivateKeyCommand(username, hostname, linuxCmd, StrictHostKeyChecking, LogLevel string, port int, connectionTimeout time.Duration) (string, []string) {
	args := []string{"-o", fmt.Sprintf("StrictHostKeyChecking=%s", StrictHostKeyChecking)}
	args = append(args, "-o", fmt.Sprintf("LogLevel=%s", LogLevel))
	args = append(args, "-i", PrivateKeyPath)
	args = append(args, fmt.Sprintf("%s@%s", username, hostname))
	args = append(args, "-p", strconv.Itoa(port))
	args = append(args, "-E", LogFilePath)
	args = append(args, fmt.Sprintf("-oConnectTimeout=%d", int(connectionTimeout.Seconds())))
	args = append(args, linuxCmd)

	return "ssh", args
}

func (s *SSHCommand) buildPasswordCommand(username, hostname, linuxCmd, StrictHostKeyChecking, LogLevel string, port int, connectionTimeout time.Duration, password string) (string, []string) {
	args := []string{}
	args = append(args, "-p", password)
	args = append(args, "ssh")
	args = append(args, "-o", fmt.Sprintf("StrictHostKeyChecking=%s", StrictHostKeyChecking))
	args = append(args, "-o", fmt.Sprintf("LogLevel=%s", LogLevel))
	args = append(args, fmt.Sprintf("%s@%s", username, hostname))
	args = append(args, "-p", strconv.Itoa(port))
	args = append(args, "-E", LogFilePath)
	args = append(args, fmt.Sprintf("-oConnectTimeout=%d", int(connectionTimeout.Seconds())))
	args = append(args, linuxCmd)

	return "sshpass", args
}

func main() {
	step.Run(&SSHCommand{})
}
