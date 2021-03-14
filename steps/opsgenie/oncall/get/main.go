package main

import (
	"encoding/json"

	envconf "github.com/caarlos0/env/v6"
	"github.com/opsgenie/opsgenie-go-sdk-v2/schedule"
	"github.com/stackpulse/public-steps/common/step"
	"github.com/stackpulse/public-steps/opsgenie/base"
)

type Args struct {
	base.Args
	ScheduleName string `env:"SCHEDULE_NAME,required"`
}

type OpsgenieOncallList struct {
	args Args
}

func (o *OpsgenieOncallList) Init() error {
	err := envconf.Parse(&o.args)
	if err != nil {
		return err
	}

	return nil
}

func (o *OpsgenieOncallList) Run() (exitCode int, output []byte, err error) {
	scheduleClient, err := schedule.NewClient(base.Config(o.args.Args))
	if err != nil {
		return step.ExitCodeFailure, nil, err
	}

	flat := false

	scheduleResult, err := scheduleClient.GetOnCalls(nil, &schedule.GetOnCallsRequest{
		Flat:                   &flat,
		ScheduleIdentifierType: schedule.Name,
		ScheduleIdentifier:     o.args.ScheduleName,
	})

	if err != nil {
		return step.ExitCodeFailure, nil, err
	}

	marshaledParticipants, err := json.Marshal(scheduleResult.OnCallParticipants)
	if err != nil {
		return step.ExitCodeFailure, nil, err
	}

	return step.ExitCodeOK, marshaledParticipants, nil
}

func main() {
	step.Run(&OpsgenieOncallList{})
}
