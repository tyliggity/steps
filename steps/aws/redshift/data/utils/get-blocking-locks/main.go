package main

import (
	"fmt"

	"github.com/stackpulse/public-steps/aws/redshift/base"
	"github.com/stackpulse/steps-sdk-go/env"
	"github.com/stackpulse/steps-sdk-go/log"
	"github.com/stackpulse/steps-sdk-go/step"
)

// See: https://github.com/awslabs/amazon-redshift-utils/blob/master/src/AdminViews/v_get_blocking_locks.sql

type redshiftAmazonUtilsGetBlockingLocks struct {
	base.RedshiftAWSRunner
	OnlyUsers       []string `env:"ONLY_USERS"`
	IgnoreUsers     []string `env:"IGNORE_USERS"`
	BlockingSeconds int      `env:"BLOCKING_SECONDS" envDefault:"60"`
}

func (r *redshiftAmazonUtilsGetBlockingLocks) Init() error {
	if err := env.Parse(r); err != nil {
		return err
	}
	log.Debug("Args: %#+v", r)

	if err := r.Validate(); err != nil {
		return err
	}

	return nil
}

func (r *redshiftAmazonUtilsGetBlockingLocks) Run() (int, []byte, error) {
	q := &query{
		onlyUsers:       r.OnlyUsers,
		ignoredUsers:    r.IgnoreUsers,
		blockingSeconds: r.BlockingSeconds,
	}
	sql, err := q.generate()
	if err != nil {
		return step.ExitCodeFailure, nil, fmt.Errorf("generate sql query: %w", err)
	}

	log.Debug("Running query: '%s'", sql)

	results, err := r.RunQuery(sql)
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
	step.Run(&redshiftAmazonUtilsGetBlockingLocks{})
}
