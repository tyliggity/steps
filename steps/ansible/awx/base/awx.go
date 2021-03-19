package base

import (
	"fmt"
	"os/exec"

	"github.com/stackpulse/steps-sdk-go/env"
	"github.com/stackpulse/steps-sdk-go/log"
	"github.com/stackpulse/steps-sdk-go/step"
)

type Args struct {
	Host     string `env:"HOST,required"`
	Username string `env:"USERNAME"`
	Password string `env:"PASSWORD"`
	Token    string `env:"TOWER_OAUTH_TOKEN"`
	Insecure bool   `env:"INSECURE" envDefault:"false"`
}

func (a Args) BaseArgs() Args {
	return a
}

func (a Args) Validate() error {
	if a.Token == "" {
		if a.Password == "" || a.Username == "" {
			return fmt.Errorf("must specify USERNAME and PASSWORD when not using token")
		}
	}
	return nil
}

type BaseArgs interface {
	BaseArgs() Args
}

type AwxCommand struct {
	Args
}

func NewAwxCommand(args BaseArgs) (*AwxCommand, error) {
	if err := env.Parse(args); err != nil {
		return nil, fmt.Errorf("parse args: %w", err)
	}
	if err := args.BaseArgs().Validate(); err != nil {
		return nil, fmt.Errorf("validate args: %w", err)
	}

	return &AwxCommand{Args: args.BaseArgs()}, nil
}

func (a *AwxCommand) Execute(extraArgs []string) (int, []byte, error) {
	args := []string{"--conf.host", a.Host, "--conf.username", a.Username}
	if a.Password != "" {
		args = append(args, "--conf.password", a.Password)
	}
	if a.Insecure {
		args = append(args, "--conf.insecure")
	}
	if env.Debug() {
		args = append(args, "--verbose")
	}

	args = append(args, extraArgs...)

	cmd := exec.Command("awx", args...)
	log.Debugln("About to run awx with args: %#v", args)

	output, err := cmd.CombinedOutput()
	if err != nil {
		if exiterr, ok := err.(*exec.ExitError); ok {
			return exiterr.ExitCode(), output, err
		}
		return step.ExitCodeFailure, output, err
	}

	return step.ExitCodeOK, output, nil
}
