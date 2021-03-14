package main

import (
	"fmt"
	"github.com/Jeffail/gabs/v2"
	envconf "github.com/caarlos0/env/v6"
	"github.com/stackpulse/steps/common/step"
	"github.com/stackpulse/steps/pagerduty/incident/base"
)

type Args struct {
	base.Args
	Id string `env:"INCIDENT_ID,required"`
}

type PagerDutyIncidentAcknowledge struct {
	args Args
}

func (p *PagerDutyIncidentAcknowledge) Init() error {
	err := envconf.Parse(&p.args)
	if err != nil {
		return err
	}

	return nil
}

func (p *PagerDutyIncidentAcknowledge) Run() (exitCode int, output []byte, err error) {
	gc := gabs.New()

	err = base.UpdateIncidentStatus(p.args.Args, p.args.Id, base.StatusAcknowledged)
	if err != nil {
		gc.Set(false, "success")

		return 1, gc.Bytes(), fmt.Errorf("failed acknowledging incident: %v", err)
	}

	gc.Set(true, "success")

	return 0, gc.Bytes(), nil
}

func main() {
	step.Run(&PagerDutyIncidentAcknowledge{})
}
