package get

import (
	"fmt"
	"strconv"

	base2 "github.com/stackpulse/steps/kubectl/base"
)

type Args struct {
	base2.Args
	Ready string `env:"READY"`
}

var NodeReadyFilterAlreadyExists = fmt.Errorf("can't add more than one filter for node ready")

func (a *Args) addReadyFilter() error {
	if a.Ready == "" {
		return nil
	}

	ready, err := strconv.ParseBool(a.Ready)
	if err != nil {
		return fmt.Errorf("can't parse READY as boolean: %w", err)
	}

	filters := a.Args.FilterEqualsParsed
	if !ready {
		filters = a.Args.FilterNotEqualsParsed
	}

	if _, ok := filters[NodeReadyJsonKey]; ok {
		return NodeReadyFilterAlreadyExists
	}

	filters[NodeReadyJsonKey] = "True"
	return nil
}

type GetNodes struct {
	Args *Args
	kctl *base2.KubectlStep
}

func NewGetNodes(args *Args) (*GetNodes, error) {
	parse := false
	if args == nil {
		parse = true
		args = &Args{}
	}
	kctl, err := base2.NewKubectlStep(args, parse)
	if err != nil {
		return nil, err
	}

	if err := args.addReadyFilter(); err != nil {
		return nil, err
	}

	return &GetNodes{
		Args: args,
		kctl: kctl,
	}, nil
}

func (n *GetNodes) Get() (output []byte, exitCode int, err error) {
	return n.kctl.Execute([]string{"get", "nodes"})
}

func (n *GetNodes) Parse(output []byte) (string, error) {
	return n.kctl.ParseOutput(output, parsingConfiguration)
}
