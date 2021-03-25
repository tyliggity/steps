package main

import (
	"github.com/Jeffail/gabs/v2"
	envconf "github.com/caarlos0/env/v6"
	"github.com/stackpulse/steps-sdk-go/step"
	"github.com/stackpulse/steps/pagerduty/incident/base"
)

type Args struct {
	base.Args
	Id string `env:"INCIDENT_ID,required"`
}
type PagerDutyIncidentResolve struct {
	args Args
}

func (p *PagerDutyIncidentResolve) Init() error {
	err := envconf.Parse(&p.args)
	if err != nil {
		return err
	}

	return nil
}

func (p *PagerDutyIncidentResolve) Run() (exitCode int, output []byte, err error) {
	gc := gabs.New()
	err = base.UpdateIncidentStatus(p.args.Args, p.args.Id, base.StatusResolved)
	if err != nil {
		gc.Set(false, "success")
		return step.ExitCodeFailure, gc.Bytes(), err
	}

	gc.Set(true, "success")

	return step.ExitCodeOK, gc.Bytes(), nil
}

func main() {
	step.Run(&PagerDutyIncidentResolve{})
}
