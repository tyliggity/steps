package main

import (
	"encoding/json"
	"fmt"

	pd "github.com/PagerDuty/go-pagerduty"
	envconf "github.com/caarlos0/env/v6"
	"github.com/stackpulse/steps/common/step"
	"github.com/stackpulse/steps/pagerduty/incident/base"
)

type Args struct {
	base.Args
	Id      string `env:"INCIDENT_ID,required"`
	Content string `env:"CONTENT,required"`
}

type PagerDutyIncidentCreateNote struct {
	args Args
}

func (p *PagerDutyIncidentCreateNote) Init() error {
	err := envconf.Parse(&p.args)
	if err != nil {
		return err
	}

	return nil
}

func (p *PagerDutyIncidentCreateNote) Run() (exitCode int, output []byte, err error) {
	client := pd.NewClient(p.args.PdToken)

	note, err := client.CreateIncidentNoteWithResponse(p.args.Id, pd.IncidentNote{
		ID:      p.args.Id,
		Content: p.args.Content,
		User:    pd.APIObject{Summary: p.args.Args.Email},
	})
	if err != nil {
		return step.ExitCodeFailure, nil, fmt.Errorf("failed to update incident: %w", err)
	}

	marshaledNote, err := json.Marshal(note)
	if err != nil {
		return step.ExitCodeFailure, nil, err
	}

	return step.ExitCodeOK, marshaledNote, nil
}

func main() {
	step.Run(&PagerDutyIncidentCreateNote{})
}
