package main

import (
	"encoding/json"
	"fmt"

	pd "github.com/PagerDuty/go-pagerduty"
	"github.com/stackpulse/public-steps/pagerduty/incident/base"
	"github.com/stackpulse/steps-sdk-go/env"
	"github.com/stackpulse/steps-sdk-go/step"
)

type pagerDutyIncidentCreate struct {
	base.Args
	ServiceName string `env:"SERVICE_NAME,required"`
	Urgency     string `env:"URGENCY,required"`
	Title       string `env:"TITLE,required"`
}

type output struct {
	Id   string `json:"id"`
	Self string `json:"self"`
	step.Outputs
}

func (p *pagerDutyIncidentCreate) Init() error {
	err := env.Parse(p)
	if err != nil {
		return err
	}

	if p.Urgency != "low" && p.Urgency != "high" {
		return fmt.Errorf("the provided urgency '%s' is invalid,  can be only 'low' or 'high'.", p.Urgency)
	}

	return nil
}

func (p *pagerDutyIncidentCreate) Run() (int, []byte, error) {
	client := pd.NewClient(p.PdToken)
	servicesResp, err := client.ListServices(pd.ListServiceOptions{})
	if err != nil {
		return step.ExitCodeFailure, nil, fmt.Errorf("list services: %w", err)
	}

	var foundService *pd.Service
	for _, s := range servicesResp.Services {
		if s.Name != p.ServiceName {
			continue
		}

		foundService = &s
		break
	}
	if foundService == nil {
		return step.ExitCodeFailure, nil, fmt.Errorf("no such service '%s': %w", p.ServiceName, err)
	}

	incident, err := client.CreateIncident(p.Email, &pd.CreateIncidentOptions{
		Type:  "incident",
		Title: p.Title,
		Service: &pd.APIReference{
			ID:   foundService.ID,
			Type: foundService.Type,
		},
		Urgency: p.Urgency,
	})
	if err != nil {
		return step.ExitCodeFailure, nil, fmt.Errorf("create incident: %w", err)
	}

	jsonOutput, err := json.Marshal(&output{Id: incident.Id, Self: incident.Self, Outputs: step.Outputs{Object: incident}})
	if err != nil {
		return step.ExitCodeFailure, nil, fmt.Errorf("marshal output: %w", err)
	}

	return step.ExitCodeOK, jsonOutput, nil
}

func main() {
	step.Run(&pagerDutyIncidentCreate{})
}
