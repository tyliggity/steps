package main

import (
	"github.com/stackpulse/steps-sdk-go/step"
	"github.com/stackpulse/steps/steps/ansible/awx/base"
)

type Args struct {
	base.Args
	ID string `env:"JOB_ID,required"`
}

type AWXJobOutput struct {
	*base.AwxCommand
	args Args
}

func (a *AWXJobOutput) Init() error {
	var args Args
	baseCommand, err := base.NewAwxCommand(&args)
	if err != nil {
		return err
	}

	a.AwxCommand = baseCommand
	a.args = args
	return nil
}

func (a *AWXJobOutput) Run() (int, []byte, error) {
	return a.Execute([]string{"jobs", "stdout", a.args.ID})
}

func main() {
	step.Run(&AWXJobOutput{})
}
