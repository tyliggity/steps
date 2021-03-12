package main

import (
	"encoding/json"
	"fmt"

	"github.com/stackpulse/public-steps/aws/redshift/base"
	"github.com/stackpulse/public-steps/common/env"
	"github.com/stackpulse/public-steps/common/log"
	"github.com/stackpulse/public-steps/common/step"
)

type redshiftCancel struct {
	base.RedshiftAWSRunner
	Pids []string `env:"PIDS,required"`
}

type output struct {
	Success []string `json:"success"`
	Failed  []string `json:"failed"`
	step.Outputs
}

func (r *redshiftCancel) Init() error {
	if err := env.Parse(r); err != nil {
		return err
	}

	log.Debug("Args: %#+v", r)

	if err := r.Validate(); err != nil {
		return err
	}

	return nil
}

func (r *redshiftCancel) Run() (int, []byte, error) {
	query, err := base.BuildPidFuncQuery(r.Pids, "pg_cancel_backend")
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
	step.Run(&redshiftCancel{})
}
