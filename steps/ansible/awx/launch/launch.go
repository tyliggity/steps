package main

import (
	"github.com/stackpulse/public-steps/common/step"
	"github.com/stackpulse/public-steps/steps/ansible/awx/base"
	"strconv"
)

type Args struct {
	base.Args
	Job       string `env:"JOB,required"`
	Wait      bool   `env:"WAIT" envDefault:"false"`
	Inventory string `env:"INVENTORY"`
	Monitor   bool   `env:"MONITOR" envDefault:"false"`
	ScmBranch string `env:"SCM_BRANCH"`
	Check     bool   `env:"JOB_CHECK" envDefault:"false"`
	Verbosity int    `env:"VERBOSITY"  envDefault:"-1"`
	ExtraVars string `env:"EXTRA_VARS"`
}

type AWXLaunch struct {
	*base.AwxCommand
	args Args
}

func (a *AWXLaunch) Init() error {
	var args Args
	baseCommand, err := base.NewAwxCommand(&args)
	if err != nil {
		return err
	}

	a.AwxCommand = baseCommand
	a.args = args
	return nil
}

func (a *AWXLaunch) getExtraArgs() []string {
	extraArgs := []string{"job_templates", "launch", a.args.Job}
	if a.args.Wait {
		extraArgs = append(extraArgs, "--wait")
	}
	if a.args.Inventory != "" {
		extraArgs = append(extraArgs, "--inventory", a.args.Inventory)
	}
	if a.args.Monitor {
		extraArgs = append(extraArgs, "--monitor")
	}
	if a.args.ScmBranch != "" {
		extraArgs = append(extraArgs, "--scm_branch", a.args.ScmBranch)
	}
	if a.args.Check {
		extraArgs = append(extraArgs, "--job_type", "check")
	}
	if a.args.Verbosity != -1 {
		extraArgs = append(extraArgs, "--verbosity", strconv.Itoa(a.args.Verbosity))
	}
	if a.args.ExtraVars != "" {
		extraArgs = append(extraArgs, "--extra_vars", a.args.ExtraVars)
	}

	return extraArgs
}

func (a *AWXLaunch) Run() (int, []byte, error) {
	extraArgs := a.getExtraArgs()
	return a.Execute(extraArgs)
}

func main() {
	step.Run(&AWXLaunch{})
}
