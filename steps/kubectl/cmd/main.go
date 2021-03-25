package main

import (
	"os"

	"github.com/stackpulse/steps-sdk-go/env"
	"github.com/stackpulse/steps-sdk-go/step"
	"github.com/stackpulse/steps/kubectl/base"
)

type KubectlCmd struct {
	kctl *base.KubectlStep
}

func (k *KubectlCmd) Init() error {
	kctl, err := base.NewKubectlStep(&base.Args{}, true)
	if err != nil {
		return err
	}
	k.kctl = kctl
	return nil
}

func (k *KubectlCmd) Run() (int, []byte, error) {
	// Setting the default formatter to raw, as the user may not produce json in this step
	env.SetFormatter(env.RawFormat, false)
	output, exitCode, err := k.kctl.Execute(os.Args[1:], base.IgnoreFieldSelector, base.IgnoreFormat, base.IgnoreNamespace)
	return exitCode, output, err
}

func main() {
	step.Run(&KubectlCmd{})
}
