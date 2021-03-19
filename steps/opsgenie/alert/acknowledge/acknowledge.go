package main

import (
	"github.com/Jeffail/gabs/v2"
	envconf "github.com/caarlos0/env/v6"
	"github.com/opsgenie/opsgenie-go-sdk-v2/alert"

	"github.com/stackpulse/public-steps/opsgenie/base"
	"github.com/stackpulse/steps-sdk-go/step"
)

type Args struct {
	base.Args
	Id string `env:"ALERT_ID,required"`
}

type OpsgenieAcknowledgeAlert struct {
	args Args
}

func (o *OpsgenieAcknowledgeAlert) Init() error {
	err := envconf.Parse(&o.args)
	if err != nil {
		return err
	}

	return nil
}

func (o *OpsgenieAcknowledgeAlert) Run() (exitCode int, output []byte, err error) {
	gc := gabs.New()

	alertClient, err := alert.NewClient(base.Config(o.args.Args))
	if err != nil {
		return step.ExitCodeFailure, nil, err
	}

	_, err = alertClient.Acknowledge(nil, &alert.AcknowledgeAlertRequest{
		IdentifierValue: o.args.Id,
	})
	if err != nil {
		return step.ExitCodeFailure, nil, err
	}

	gc.Set(true, "success")

	return step.ExitCodeOK, gc.Bytes(), nil
}

func main() {
	step.Run(&OpsgenieAcknowledgeAlert{})
}
