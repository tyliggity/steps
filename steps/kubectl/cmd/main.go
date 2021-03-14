package main

import (
	"github.com/stackpulse/public-steps/common/env"
	"github.com/stackpulse/public-steps/common/step"
	"github.com/stackpulse/public-steps/kubectl/base"
	"os"
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
