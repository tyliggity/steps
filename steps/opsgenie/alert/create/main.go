package main

import (
	"fmt"
	"reflect"

	"github.com/Jeffail/gabs/v2"
	envconf "github.com/caarlos0/env/v6"
	"github.com/opsgenie/opsgenie-go-sdk-v2/alert"
	"github.com/stackpulse/public-steps/common/env"
	"github.com/stackpulse/public-steps/common/step"
	"github.com/stackpulse/public-steps/opsgenie/base"
)

type Args struct {
	base.Args
	Message       string            `env:"MESSAGE,required"`
	Description   string            `env:"DESCRIPTION,required"`
	Priority      string            `env:"PRIORITY" envDefault:"P3"`
	Details       map[string]string `env:"DETAILS"`
	ResponderType string            `env:"RESPONDER_TYPE"`
	ResponderName string            `env:"RESPONDER_NAME"`
}

type OpsgenieCreateAlert struct {
	args Args
}

func validateResponder(args Args) error {
	switch alert.ResponderType(args.ResponderType) {
	case alert.TeamResponder, alert.EscalationResponder, alert.UserResponder, alert.ScheduleResponder:
		if args.ResponderName != "" {
			return nil
		}
		return fmt.Errorf("responder name is empty")
	case "":
		return nil
	}

	return fmt.Errorf("invalid responder type, should be one of: 'user', 'team', 'escalation', 'schedule'")

}

func (o *OpsgenieCreateAlert) Init() error {
	err := envconf.ParseWithFuncs(&o.args, map[reflect.Type]envconf.ParserFunc{
		reflect.TypeOf(map[string]string{}): env.ParseKeyValueEnv,
	})

	if err != nil {
		return err
	}

	err = validateResponder(o.args)
	if err != nil {
		return err
	}

	return alert.ValidatePriority(alert.Priority(o.args.Priority))
}

func (o *OpsgenieCreateAlert) Run() (exitCode int, output []byte, err error) {
	gc := gabs.New()

	alertClient, err := alert.NewClient(base.Config(o.args.Args))
	if err != nil {
		return step.ExitCodeFailure, nil, err
	}

	var responders []alert.Responder

	if o.args.ResponderType != "" {
		responders = []alert.Responder{
			{
				Type: alert.ResponderType(o.args.ResponderType),
				Name: o.args.ResponderName,
			},
		}
	}

	_, err = alertClient.Create(nil, &alert.CreateAlertRequest{
		Message:     o.args.Message,
		Description: o.args.Description,
		Priority:    alert.Priority(o.args.Priority),
		Responders:  responders,
		Details:     o.args.Details,
	})

	if err != nil {
		return step.ExitCodeFailure, nil, err
	}

	gc.Set(true, "success")

	return step.ExitCodeOK, gc.Bytes(), nil
}

func main() {
	step.Run(&OpsgenieCreateAlert{})
}
