package main

import (
	"encoding/json"
	"fmt"
	"github.com/stackpulse/steps/prometheus/base"
	rangeBase "github.com/stackpulse/steps/prometheus/query/range/base"
)

type Args struct {
	rangeBase.Args
	rangeBase.BaseArgs
	rangeBase.DurationOrExplicitArgs
}

func run() error {
	args := &Args{}
	if err := rangeBase.ParseArgs(args); err != nil {
		return err
	}

	startTime, err := rangeBase.ParseTime(args.StartTime)
	if err != nil {
		return fmt.Errorf("failed parsing start time: %w", err)
	}

	rangeOutput, err := rangeBase.DurationOrRaw(&args.Args, &args.DurationOrExplicitArgs, args.Query, startTime, args.Step)
	if err != nil {
		return err
	}

	output, err := json.Marshal(rangeOutput)
	if err != nil {
		return fmt.Errorf("can't marshal range output: %w", err)
	}

	fmt.Print(string(output))
	return nil
}

func main() {
	if err := run(); err != nil {
		base.ErrorAndExit(err.Error())
	}
}
