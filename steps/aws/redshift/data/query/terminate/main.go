package main

import (
	"encoding/json"
	"fmt"

	"github.com/stackpulse/public-steps/redshift/base"

	"github.com/stackpulse/steps-sdk-go/env"
	"github.com/stackpulse/steps-sdk-go/log"
	"github.com/stackpulse/steps-sdk-go/step"
)

type redshiftTerminateSession struct {
	base.RedshiftAWSRunner
	Pids []string `env:"PIDS,required"`
}

type output struct {
	Success []string `json:"success"`
	Failed  []string `json:"failed"`
	step.Outputs
}

func (r *redshiftTerminateSession) Init() error {
	if err := env.Parse(r); err != nil {
		return err
	}

	log.Debug("Args: %#+v", r)

	if err := r.Validate(); err != nil {
		return err
	}

	return nil
}

func (r *redshiftTerminateSession) Run() (int, []byte, error) {
	query, err := base.BuildPidFuncQuery(r.Pids, "pg_terminate_backend")
	if err != nil {
		return step.ExitCodeFailure, nil, fmt.Errorf("build pids cancel query for '%s': %w", r.Pids, err)
	}

	log.Debug("Running query: '%s'", query)
	results, err := r.RunQuery(query)
	if err != nil {
		return step.ExitCodeFailure, nil, fmt.Errorf("run query: %w", err)
	}

	success, failed, err := base.ParsePidQueryResults(results)
	if err != nil {
		return step.ExitCodeFailure, nil, fmt.Errorf("parse pid query results: %w", err)
	}

	jsonOutput, err := json.Marshal(&output{Success: success, Failed: failed})
	if err != nil {
		return step.ExitCodeFailure, nil, fmt.Errorf("marshal output: %w", err)
	}

	return step.ExitCodeOK, jsonOutput, nil
}

func main() {
	step.Run(&redshiftTerminateSession{})
}
