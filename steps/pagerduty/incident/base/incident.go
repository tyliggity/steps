package base

import (
	"fmt"

	pd "github.com/PagerDuty/go-pagerduty"
	"github.com/stackpulse/public-steps/pagerduty/base"
)

type IncidentStatus string

const (
	StatusResolved     = "resolved"
	StatusAcknowledged = "acknowledged"
)

type Args struct {
	base.Args
	Email string `env:"PAGERDUTY_USER_EMAIL,required"`
}

func UpdateIncidentStatus(args Args, incidentID string, status IncidentStatus) error {
	client := pd.NewClient(args.Args.PdToken)

	incidentOpts := pd.ManageIncidentsOptions{
		ID:     incidentID,
		Type:   "incident",
		Status: string(status),
	}

	_, err := client.ManageIncidents(args.Email, []pd.ManageIncidentsOptions{incidentOpts})
	if err != nil {
		return fmt.Errorf("update incident: %w", err)
	}

	return nil
}
