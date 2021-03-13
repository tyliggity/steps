package main

import (
	"encoding/json"
	"fmt"

	pd "github.com/PagerDuty/go-pagerduty"

	"github.com/stackpulse/public-steps/common/env"
	"github.com/stackpulse/public-steps/common/log"
	"github.com/stackpulse/public-steps/common/step"
	pgBase "github.com/stackpulse/public-steps/pagerduty/base"
	"github.com/stackpulse/public-steps/pagerduty/incident/base"
)

type pagerDutyIncidentList struct {
	base.Args
	ServiceIDs   []string `env:"SERVICE_IDS" envSeparator:","`
	ServiceNames []string `env:"SERVICE_NAMES" envSeparator:","`
	TimeZone     string   `env:"TIME_ZONE"`
	Since        string   `env:"SINCE"`
	Until        string   `env:"UNTIL"`
	Limit        int      `env:"LIMIT"`
}

type output struct {
	Total     int           `json:"total"`
	Incidents []pd.Incident `json:"incidents"`
}

func (p *pagerDutyIncidentList) Init() error {
	err := env.Parse(p)
	if err != nil {
		return err
	}

	return nil
}

func (p *pagerDutyIncidentList) buildListOpts() (pd.ListIncidentsOptions, error) {
	incidentOpts := pd.ListIncidentsOptions{}

	if len(p.ServiceNames) > 0 {
		idsByNames, err := pgBase.ServiceIdsByNames(p.Args.Args, p.ServiceNames)
		if err != nil {
			return incidentOpts, fmt.Errorf("service ids by names: %w", err)
		}
		p.ServiceIDs = append(p.ServiceIDs, idsByNames...)
		log.Debug("Listing incidents with service names %v", p.ServiceNames)
	}

	if len(p.ServiceIDs) > 0 {
		incidentOpts.ServiceIDs = p.ServiceIDs
		log.Debug("Listing incidents with service ids %v", p.ServiceIDs)
	}

	if p.Since != "" {
		incidentOpts.Since = p.Since
		log.Debug("Listing incidents since %v", p.Since)
	}

	if p.Until != "" {
		incidentOpts.Until = p.Until
		log.Debug("Listing incidents until %v", p.Until)
	}

	if p.TimeZone != "" {
		incidentOpts.TimeZone = p.TimeZone
		log.Debug("Listing incidents with time zone %v", p.TimeZone)
	}

	if p.Limit > 0 {
		incidentOpts.Limit = uint(p.Limit)
		log.Debug("Listing incidents with limit %v", p.Limit)
	}

	return incidentOpts, nil
}

func (p *pagerDutyIncidentList) Run() (int, []byte, error) {
	client := pd.NewClient(p.PdToken)

	opts, err := p.buildListOpts()
	if err != nil {
		return step.ExitCodeFailure, nil, fmt.Errorf("build list options: %w", err)
	}

	var incidents []pd.Incident
	resp, err := client.ListIncidents(opts)
	if err != nil {
		return step.ExitCodeFailure, nil, fmt.Errorf("first list incidents: %w", err)
	}

	incidents = append(incidents, resp.Incidents...)

	if p.Limit == 0 {
		for resp.More {
			opts.Offset += resp.Limit
			nextPage, err := client.ListIncidents(opts)
			if err != nil {
				return step.ExitCodeFailure, nil, fmt.Errorf("list incidents (offest: %d): %w", opts.Offset, err)
			}
			incidents = append(incidents, nextPage.Incidents...)
			resp.More = nextPage.More
		}
	}

	jsonOutput, err := json.Marshal(&output{Total: len(incidents), Incidents: incidents})
	if err != nil {
		return step.ExitCodeFailure, nil, fmt.Errorf("marshal output: %w", err)
	}

	return step.ExitCodeOK, jsonOutput, nil
}

func main() {
	step.Run(&pagerDutyIncidentList{})
}
