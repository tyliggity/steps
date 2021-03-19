package get

import (
	"fmt"

	"github.com/Jeffail/gabs/v2"
	base2 "github.com/stackpulse/public-steps/kubectl/base"
	"github.com/stackpulse/steps-sdk-go/env"
)

type Args struct {
	base2.Args
	Short      bool `env:"SHORT" envDefault:"false"`
	ClientOnly bool `env:"CLIENT_ONLY" envDefault:"false"`
}

type GetVersion struct {
	Args *Args
	kctl *base2.KubectlStep
}

func NewGetVersion(args *Args) (*GetVersion, error) {
	parse := false
	if args == nil {
		parse = true
		args = &Args{}
	}
	kctl, err := base2.NewKubectlStep(args, parse)
	if err != nil {
		return nil, err
	}

	return &GetVersion{
		Args: args,
		kctl: kctl,
	}, nil
}

func (v *GetVersion) Get() (output []byte, exitCode int, err error) {
	cmdArgs := []string{
		"version", fmt.Sprintf("--short=%v", v.Args.Short), fmt.Sprintf("--client=%v", v.Args.ClientOnly),
	}
	return v.kctl.Execute(cmdArgs)
}

func (v *GetVersion) Parse(output []byte) (string, error) {
	if !env.FormatterIs(env.JsonFormat) {
		return string(output), nil
	}

	gc, err := gabs.ParseJSON(output)
	if err != nil {
		return "", fmt.Errorf("can't parse output: %w", err)
	}

	retGc := gc
	if v.Args.Short {
		retGc = gabs.New()

		if !gc.Exists("clientVersion") {
			return "", fmt.Errorf("can't find clientVersion key in output")
		}
		retGc.Set(gc.S("clientVersion", "gitVersion").Data(), "clientVersion", "version")

		if gc.Exists("serverVersion") {
			retGc.Set(gc.S("serverVersion", "gitVersion").Data(), "serverVersion", "version")
		}
	}

	if v.Args.Pretty {
		return retGc.StringIndent("", "  "), nil
	}
	return retGc.String(), nil
}
