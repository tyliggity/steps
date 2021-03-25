package get

import (
	base2 "github.com/stackpulse/steps/kubectl/base"
)

type Args struct {
	base2.Args
	StatefulsetNames []string `env:"STATEFULSET_NAMES"`
}

func (a *Args) addFilters() error {
	return nil
}

type Statefulset struct {
	Args *Args
	kctl *base2.KubectlStep
}

func NewGetStatefulset(args *Args) (*Statefulset, error) {
	parse := false
	if args == nil {
		parse = true
		args = &Args{}
	}
	kctl, err := base2.NewKubectlStep(args, parse)
	if err != nil {
		return nil, err
	}

	if err := args.addFilters(); err != nil {
		return nil, err
	}

	return &Statefulset{
		Args: args,
		kctl: kctl,
	}, nil
}

func (n *Statefulset) Get() (output []byte, exitCode int, err error) {
	cmdArgs := []string{"get", "statefulsets"}
	cmdArgs = append(cmdArgs, n.Args.StatefulsetNames...)
	return n.kctl.Execute(cmdArgs)
}

func (n *Statefulset) Parse(output []byte) (string, error) {
	return n.kctl.ParseOutput(output, parsingConfiguration)
}
