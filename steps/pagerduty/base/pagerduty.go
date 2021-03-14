package base

import (
	"fmt"

	pd "github.com/PagerDuty/go-pagerduty"
)

type Args struct {
	PdToken string `env:"PAGERDUTY_TOKEN,required"`
}

func (a *Args) PagerDutyArgs() Args {
	return *a
}

type PagerDutyArgs interface {
	PagerDutyArgs() Args
}

func ServiceIdsByNames(args Args, serviceNames []string) ([]string, error) {
	pdClient := pd.NewClient(args.PdToken)
	servicesResp, err := pdClient.ListServices(pd.ListServiceOptions{})
	if err != nil {
		return nil, fmt.Errorf("list services: %w", err)
	}

	var serviceIDs []string
	for _, service := range servicesResp.Services {
		if Contains(serviceNames, service.Name) {
			serviceIDs = append(serviceIDs, service.ID)
		}
		if len(serviceIDs) >= len(serviceNames) {
			break
		}
	}
	return serviceIDs, nil
}

func Contains(a []string, x string) bool {
	for _, n := range a {
		if x == n {
			return true
		}
	}
	return false
}
