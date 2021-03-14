package main

import (
	"github.com/stackpulse/public-steps/common/step"
	"github.com/stackpulse/public-steps/steps/ansible/awx/base"
)

type Args struct {
	base.Args
	ID string `env:"JOB_ID,required"`
}

type AWXJobEvents struct {
	*base.AwxCommand
	args Args
}

func (a *AWXJobEvents) Init() error {
	var args Args
	baseCommand, err := base.NewAwxCommand(&args)
	if err != nil {
		return err
	}

	a.AwxCommand = baseCommand
	a.args = args
	return nil
}

func (a *AWXJobEvents) Run() (int, []byte, error) {
	return a.Execute([]string{"job_events", "get", a.args.ID})
}

func main() {
	step.Run(&AWXJobEvents{})
}
