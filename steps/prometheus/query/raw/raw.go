package main

import (
	"encoding/json"
	"fmt"
	envconf "github.com/caarlos0/env/v6"
	"github.com/stackpulse/public-steps/prometheus/base"
	queryBase "github.com/stackpulse/public-steps/prometheus/query/base"
	"time"
)

type Args struct {
	queryBase.Args
	Query string `env:"QUERY,required"`
	Time  string `env:"TIME"`
}

func run() error {
	args := &Args{}
	if err := envconf.Parse(args); err != nil {
		return err
	}

	evaluationTime := time.Time{}
	if args.Time != "" {
		t, err := time.Parse(time.RFC3339, args.Time)
		if err != nil {
			return fmt.Errorf("evaluation time parse failed: %w", err)
		}
		evaluationTime = t
	}

	res, err := queryBase.Query(&args.Args, args.Query, evaluationTime)
	if err != nil {
		return fmt.Errorf("prometheus query failed: %w", err)
	}

	output, err := json.Marshal(res)
	if err != nil {
		return fmt.Errorf("prometheus output marshaling failed: %w", err)
	}

	fmt.Print(string(output))
	return nil
}

func main() {
	if err := run(); err != nil {
		base.ErrorAndExit(err.Error())
	}
}
