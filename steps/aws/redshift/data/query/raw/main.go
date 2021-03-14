package main

import (
	"fmt"

	"github.com/stackpulse/public-steps/common/env"
	"github.com/stackpulse/public-steps/common/log"
	"github.com/stackpulse/public-steps/common/step"
	"github.com/stackpulse/public-steps/redshift/base"
)

type redshiftDataExecute struct {
	base.RedshiftAWSRunner
	Query string `env:"QUERY,required"`
}

func (r *redshiftDataExecute) Init() error {
	if err := env.Parse(r); err != nil {
		return err
	}

	log.Debug("Args: %#+v", r)

	if err := r.Validate(); err != nil {
		return err
	}

	return nil
}

func (r *redshiftDataExecute) Run() (int, []byte, error) {
	results, err := r.RunQuery(r.Query)
	if err != nil {
		return step.ExitCodeFailure, nil, fmt.Errorf("run query: %w", err)
	}

	jsonOutput, err := results.ResultsOutput()
	if err != nil {
		return step.ExitCodeFailure, nil, fmt.Errorf("results output: %w", err)
	}

	return step.ExitCodeOK, jsonOutput, nil
}

func main() {
	step.Run(&redshiftDataExecute{})
}
