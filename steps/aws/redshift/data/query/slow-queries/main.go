package main

import (
	"fmt"

	"github.com/stackpulse/steps/redshift/base"

	"github.com/stackpulse/steps-sdk-go/env"
	"github.com/stackpulse/steps-sdk-go/log"
	"github.com/stackpulse/steps-sdk-go/step"
)

type redshiftDataSlowQueries struct {
	base.RedshiftAWSRunner
	Limit int `env:"LIMIT"`
}

func (r *redshiftDataSlowQueries) Init() error {
	if err := env.Parse(r); err != nil {
		return err
	}

	log.Debug("Args: %#+v", r)

	if err := r.Validate(); err != nil {
		return err
	}

	return nil
}

func (r *redshiftDataSlowQueries) Run() (int, []byte, error) {
	query, err := generate(r.Limit)
	if err != nil {
		return step.ExitCodeFailure, nil, fmt.Errorf("generate query: %w", err)
	}

	results, err := r.RunQuery(query)
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
	step.Run(&redshiftDataSlowQueries{})
}
