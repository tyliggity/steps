package base

import (
	"context"
	"fmt"

	"github.com/bndr/gojenkins"
	"github.com/stackpulse/public-steps/common/env"
)

type Args struct {
	Username string `env:"USERNAME,required"`
	Password string `env:"PASSWORD,required"`
	Host     string `env:"HOST,required"`
}

func (a Args) BaseArgs() Args {
	return a
}

type BaseArgs interface {
	BaseArgs() Args
}

type JenkinsBaseCommand struct {
	Args
	client *gojenkins.Jenkins
}

func NewBaseJenkinsCommand(args BaseArgs) (*JenkinsBaseCommand, error) {
	if err := env.Parse(args); err != nil {
		return nil, fmt.Errorf("parse args: %w", err)
	}

	baseArgs := args.BaseArgs()
	jenkins, err := gojenkins.CreateJenkins(nil, baseArgs.Host, baseArgs.Username, baseArgs.Password).Init(context.Background())
	if err != nil {
		return nil, fmt.Errorf("init jenkins: %w", err)
	}
	return &JenkinsBaseCommand{Args: baseArgs, client: jenkins}, nil
}

func (j *JenkinsBaseCommand) Client() *gojenkins.Jenkins {
	return j.client
}
