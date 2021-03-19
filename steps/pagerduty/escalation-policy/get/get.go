package main

import (
	"encoding/json"
	"fmt"

	pd "github.com/PagerDuty/go-pagerduty"

	"github.com/stackpulse/public-steps/pagerduty/base"
	"github.com/stackpulse/steps-sdk-go/env"
	"github.com/stackpulse/steps-sdk-go/step"
)

type pagerDutyEscalationPolicyGet struct {
	base.Args
	ID string `env:"ESCALATION_POLICY_ID,required"`
}

type output struct {
	Id          string `json:"id"`
	Summary     string `json:"summary"`
	Description string `json:"description"`
	Self        string `json:"self"`
	step.Outputs
}

func (p *pagerDutyEscalationPolicyGet) Init() error {
	err := env.Parse(p)
	if err != nil {
		return fmt.Errorf("config parse: %w", err)
	}

	return nil
}

func (p *pagerDutyEscalationPolicyGet) Run() (int, []byte, error) {
	client := pd.NewClient(p.PdToken)
	resp, err := client.GetEscalationPolicy(p.ID, &pd.GetEscalationPolicyOptions{})
	if err != nil {
		return step.ExitCodeFailure, nil, fmt.Errorf("get escalation policy: %w", err)
	}

	jsonOutput, err := json.Marshal(&output{
		Id:          resp.ID,
		Summary:     resp.Summary,
		Description: resp.Description,
		Self:        resp.Self,
		Outputs:     step.Outputs{Object: resp},
	})
	if err != nil {
		return step.ExitCodeFailure, nil, fmt.Errorf("marshal output: %w", err)
	}

	return step.ExitCodeOK, jsonOutput, nil
}

func main() {
	step.Run(&pagerDutyEscalationPolicyGet{})
}
