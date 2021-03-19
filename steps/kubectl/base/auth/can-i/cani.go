package can_i

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/Jeffail/gabs"
	"github.com/google/shlex"
	"github.com/stackpulse/public-steps/kubectl/base"
	"github.com/stackpulse/steps-sdk-go/env"
)

type Args struct {
	base.Args
	Resource string `env:"RESOURCE,required"`
}

type CanI struct {
	Args *Args
	kctl *base.KubectlStep
}

func NewCanI(args *Args) (*CanI, error) {
	parse := false
	if args == nil {
		parse = true
		args = &Args{}
	}

	kctl, err := base.NewKubectlStep(args, parse)
	if err != nil {
		return nil, err
	}

	return &CanI{
		Args: args,
		kctl: kctl,
	}, nil
}

func (c *CanI) Run() (output []byte, exitCode int, err error) {
	extraArgs, err := shlex.Split(c.Args.Resource)
	if err != nil {
		return nil, 1, fmt.Errorf("error splitting resource using shlex: %w", err)
	}
	cmdArgs := []string{"auth", "can-i"}
	cmdArgs = append(cmdArgs, extraArgs...)

	output, exitCode, err = c.kctl.Execute(cmdArgs, base.IgnoreFormat, base.IgnoreFieldSelector)
	if err != nil {
		// If the answer is no, so an error will be returned, but in this case, it's not a real error :)
		if exitCode == 1 && output != nil && strings.HasPrefix(string(output), "no") {
			exitCode = 0
			err = nil
		}
	}

	return output, exitCode, err
}

func (c *CanI) Parse(output []byte) (string, error) {
	outputStr := string(output)

	if !env.FormatterIs(env.JsonFormat) {
		return outputStr, nil
	}

	gc := gabs.New()

	cani := strings.HasPrefix(outputStr, "yes")
	gc.Set(cani, "cani")

	if !cani {
		r := regexp.MustCompile("no - (.+)")
		found := r.FindStringSubmatch(outputStr)
		if len(found) > 1 {
			gc.Set(found[1], "reason")
		}
	}

	if c.Args.Pretty {
		return gc.StringIndent("", "  "), nil
	}
	return gc.String(), nil
}
