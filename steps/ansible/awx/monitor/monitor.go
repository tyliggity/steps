package main

import (
	"github.com/stackpulse/public-steps/steps/ansible/awx/base"
	"github.com/stackpulse/steps-sdk-go/step"
)

type Args struct {
	base.Args
	ID string `env:"JOB_ID,required"`
}

type AWXMonitor struct {
	*base.AwxCommand
	args Args
}

func (a *AWXMonitor) Init() error {
	var args Args
	baseCommand, err := base.NewAwxCommand(&args)
	if err != nil {
		return err
	}

	a.AwxCommand = baseCommand
	a.args = args
	return nil
}

func (a *AWXMonitor) Run() (int, []byte, error) {
	return a.Execute([]string{"jobs", "monitor", a.args.ID})
}

func main() {
	step.Run(&AWXMonitor{})
}
