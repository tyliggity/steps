package main

import (
	"github.com/stackpulse/public-steps/steps/ansible/awx/base"
	"github.com/stackpulse/steps-sdk-go/step"
)

type Args struct {
	base.Args
}

type AWXHosts struct {
	*base.AwxCommand
	args Args
}

func (a *AWXHosts) Init() error {
	var args Args
	baseCommand, err := base.NewAwxCommand(&args)
	if err != nil {
		return err
	}

	a.AwxCommand = baseCommand
	a.args = args
	return nil
}

func (a *AWXHosts) Run() (int, []byte, error) {
	return a.Execute([]string{"hosts", "list"})
}

func main() {
	step.Run(&AWXHosts{})
}
