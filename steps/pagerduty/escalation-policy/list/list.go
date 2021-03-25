package main

import (
	"encoding/json"
	"fmt"

	pd "github.com/PagerDuty/go-pagerduty"

	"github.com/stackpulse/steps-sdk-go/env"
	"github.com/stackpulse/steps-sdk-go/step"
	"github.com/stackpulse/steps/pagerduty/base"
)

type pagerDutyEscalationPoliciesList struct {
	base.Args
}

type output struct {
	Policies []pd.EscalationPolicy `json:"policies"`
}

func (p *pagerDutyEscalationPoliciesList) Init() error {
	err := env.Parse(p)
	if err != nil {
		return fmt.Errorf("config parse: %w", err)
	}

	return nil
}

func (p *pagerDutyEscalationPoliciesList) Run() (int, []byte, error) {
	client := pd.NewClient(p.PdToken)
	opts := pd.ListEscalationPoliciesOptions{}
	var policies []pd.EscalationPolicy
	resp, err := client.ListEscalationPolicies(opts)
	if err != nil {
		return step.ExitCodeFailure, nil, fmt.Errorf("first list escalation policies: %w", err)
	}
	policies = resp.EscalationPolicies

	for resp.More {
		opts.Offset += resp.Limit
		nextPage, err := client.ListEscalationPolicies(opts)
		if err != nil {
			return step.ExitCodeFailure, nil, fmt.Errorf("list escalation policies (offset %d): %w", opts.Offset, err)
		}
		policies = append(policies, nextPage.EscalationPolicies...)
		resp.More = nextPage.More
	}

	jsonOutput, err := json.Marshal(&output{Policies: policies})
	if err != nil {
		return step.ExitCodeFailure, nil, fmt.Errorf("marshal output: %w", err)
	}

	return step.ExitCodeOK, jsonOutput, nil
}

func main() {
	step.Run(&pagerDutyEscalationPoliciesList{})
}
