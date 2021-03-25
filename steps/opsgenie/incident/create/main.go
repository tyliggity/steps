package main

import (
	"fmt"

	"github.com/Jeffail/gabs/v2"
	envconf "github.com/caarlos0/env/v6"
	"github.com/opsgenie/opsgenie-go-sdk-v2/incident"
	"github.com/opsgenie/opsgenie-go-sdk-v2/service"
	"github.com/stackpulse/steps-sdk-go/step"
	"github.com/stackpulse/steps/opsgenie/base"
)

type Args struct {
	base.Args
	ServiceName string `env:"SERVICE_NAME,required"`
	Message     string `env:"MESSAGE,required"`
	Description string `env:"DESCRIPTION,required"`
	Priority    string `env:"PRIORITY" envDefault:"P3"`
}

type OpsgenieIncidentCreate struct {
	args Args
}

func (o *OpsgenieIncidentCreate) Init() error {
	err := envconf.Parse(&o.args)
	if err != nil {
		return err
	}

	return incident.ValidatePriority(incident.Priority(o.args.Priority))
}

func (o *OpsgenieIncidentCreate) getServiceByName(serviceName string) (*service.Service, error) {
	serviceClient, err := service.NewClient(base.Config(o.args.Args))

	if err != nil {
		return nil, err
	}

	currentLimit := 100 // 100 is the allowed maximum
	currentOffset := 0

	for {
		listResult, err := serviceClient.List(nil, &service.ListRequest{
			Limit:  currentLimit,
			Offset: currentOffset,
		})

		if err != nil {
			return nil, err
		}

		for _, s := range listResult.Services {

			if s.Name == serviceName {
				return &s, nil
			}
		}

		if len(listResult.Services) < currentLimit {
			break
		}

		currentOffset += currentLimit
	}

	return nil, fmt.Errorf("no such service: '%s'", serviceName)
}

func (o *OpsgenieIncidentCreate) Run() (exitCode int, output []byte, err error) {
	gc := gabs.New()

	incidentClient, err := incident.NewClient(base.Config(o.args.Args))
	if err != nil {
		return step.ExitCodeFailure, nil, err
	}

	opsgenieService, err := o.getServiceByName(o.args.ServiceName)
	if err != nil {
		return step.ExitCodeFailure, nil, err
	}

	_, err = incidentClient.Create(nil, &incident.CreateRequest{
		Message:     o.args.Message,
		Description: o.args.Description,
		Priority:    incident.Priority(o.args.Priority),
		ServiceId:   opsgenieService.Id,
	})

	if err != nil {
		return step.ExitCodeFailure, nil, err
	}

	gc.Set(true, "success")

	return step.ExitCodeOK, gc.Bytes(), nil
}

func main() {
	step.Run(&OpsgenieIncidentCreate{})
}
