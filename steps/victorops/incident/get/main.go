package main

import (
	"encoding/json"

	"github.com/stackpulse/public-steps/victorops/base"
	"github.com/stackpulse/steps-sdk-go/env"
	"github.com/stackpulse/steps-sdk-go/step"

	"github.com/victorops/go-victorops/victorops"
)

type IncidentGet struct {
	base.Args
	IncidentNumber int `env:"INCIDENT_NUMBER,required"`
}

type getRequest struct {
	IncidentNumbers []string `json:"incidentNames"`
	UserName        string   `json:"userName"`
	ResolveMessage  string   `json:"message"`
}

type output struct {
	CurrentPhase string `json:"current_phase"`
	AlertCount   int    `json:"alert_count"`
	EntityID     string `json:"entity_id"`
	step.Outputs
}

func (s *IncidentGet) Init() error {
	err := env.Parse(s)
	if err != nil {
		return err
	}

	return nil
}

func (s *IncidentGet) Run() (int, []byte, error) {
	victoropsClient := victorops.NewClient(s.ApiID, s.ApiKey, "https://api.victorops.com")

	incident, details, err := victoropsClient.GetIncident(s.IncidentNumber)
	if err != nil {
		return step.ExitCodeFailure, nil, err
	}

	//get raw response body object
	var respBody map[string]interface{}
	err = json.Unmarshal([]byte(details.ResponseBody), &respBody)
	if err != nil {
		return step.ExitCodeFailure, nil, err
	}

	//prepare outputs
	stepOutput := output{
		CurrentPhase: incident.CurrentPhase,
		AlertCount:   incident.AlertCount,
		EntityID:     incident.EntityID,
		Outputs:      step.Outputs{Object: respBody},
	}
	jsonOutput, err := json.Marshal(stepOutput)
	if err != nil {
		return step.ExitCodeFailure, nil, err
	}

	return step.ExitCodeOK, jsonOutput, nil
}

func main() {
	step.Run(&IncidentGet{})
}
