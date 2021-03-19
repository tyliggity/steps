package main

import (
	"github.com/Jeffail/gabs/v2"
	envconf "github.com/caarlos0/env/v6"
	"github.com/opsgenie/opsgenie-go-sdk-v2/incident"

	"github.com/stackpulse/public-steps/opsgenie/base"
	"github.com/stackpulse/steps-sdk-go/step"
)

type Args struct {
	base.Args
	Id   string `env:"INCIDENT_ID,required"`
	Note string `env:"NOTE,required"`
}

type OpsgenieIncidentAddNote struct {
	args Args
}

func (o *OpsgenieIncidentAddNote) Init() error {
	err := envconf.Parse(&o.args)
	if err != nil {
		return err
	}

	return nil
}

func (o *OpsgenieIncidentAddNote) Run() (exitCode int, output []byte, err error) {
	gc := gabs.New()

	incidentClient, err := incident.NewClient(base.Config(o.args.Args))
	if err != nil {
		return step.ExitCodeFailure, nil, err
	}

	_, err = incidentClient.AddNote(nil, &incident.AddNoteRequest{
		Id:         o.args.Id,
		Identifier: incident.Id,
		Note:       o.args.Note,
	})
	if err != nil {
		return step.ExitCodeFailure, nil, err
	}

	gc.Set(true, "success")

	return step.ExitCodeOK, gc.Bytes(), nil
}

func main() {
	step.Run(&OpsgenieIncidentAddNote{})
}
