package get

import (
	base2 "github.com/stackpulse/steps/kubectl/base"
)

type Args struct {
	base2.Args
	ConfigmapNames []string `env:"CONFIGMAP_NAMES"`
}

type GetConfigmap struct {
	Args *Args
	kctl *base2.KubectlStep
}

func NewGetConfigmap(args *Args) (*GetConfigmap, error) {
	parse := false
	if args == nil {
		parse = true
		args = &Args{}
	}
	kctl, err := base2.NewKubectlStep(args, parse)
	if err != nil {
		return nil, err
	}

	return &GetConfigmap{
		Args: args,
		kctl: kctl,
	}, nil
}

func (l *GetConfigmap) Get() (output []byte, exitCode int, err error) {
	cmdArgs := []string{"get", "configmap"}
	cmdArgs = append(cmdArgs, l.Args.ConfigmapNames...)
	return l.kctl.Execute(cmdArgs)
}

func (l *GetConfigmap) Parse(output []byte) (string, error) {
	return l.kctl.ParseOutput(output, nil)
}
