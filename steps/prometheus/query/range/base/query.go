package query

import (
	"fmt"
	"time"

	envconf "github.com/caarlos0/env/v6"
	"github.com/prometheus/common/model"
	queryBase "github.com/stackpulse/steps/prometheus/query/base"
)

type BaseArgs struct {
	Query     string        `env:"QUERY,required"`
	StartTime string        `env:"START,required"`
	Step      time.Duration `env:"STEP" envDefault:"5h"`
}

type DurationArgs struct {
	Since time.Duration `env:"SINCE,required"`
}

type DurationOrExplicitArgs struct {
	Since   time.Duration `env:"SINCE"`
	EndTime string        `env:"END"`
}

type Args struct {
	queryBase.Args
}

func ParseArgs(args interface{}) error {
	if err := envconf.Parse(args); err != nil {
		return err
	}
	if durationOrExplicit, ok := args.(*DurationOrExplicitArgs); ok {
		if durationOrExplicit.Since == 0 && durationOrExplicit.EndTime == "" {
			return fmt.Errorf("at least one of 'SINCE' or 'END' args cannot be empty")
		}
	}
	return nil
}

func MatrixDurationOrRaw(args *Args, options *DurationOrExplicitArgs, query string, startTime time.Time, step time.Duration) (model.Matrix, error) {
	val, err := DurationOrRaw(args, options, query, startTime, step)
	if err != nil {
		return nil, err
	}
	metrics, ok := val.(model.Matrix)
	if !ok {
		return nil, fmt.Errorf("%T output format is not supported", val)
	}
	return metrics, nil
}

func DurationOrRaw(args *Args, options *DurationOrExplicitArgs, query string, startTime time.Time, step time.Duration) (model.Value, error) {
	var rangeOutput model.Value
	var err error
	if options.EndTime != "" {
		endTime, err := ParseTime(options.EndTime)
		if err != nil {
			return nil, fmt.Errorf("failed parsing end time: %w", err)
		}

		rangeOutput, err = Range(args, query, startTime, endTime, step)
		if err != nil {
			return nil, err
		}
	} else {
		rangeOutput, err = RangeDuration(args, query, startTime, options.Since, step)
	}
	return rangeOutput, err
}

func MatrixDuration(args *Args, query string, baseTime time.Time, since time.Duration, step time.Duration) (model.Matrix, error) {
	val, err := RangeDuration(args, query, baseTime, since, step)
	if err != nil {
		return nil, err
	}
	metrics, ok := val.(model.Matrix)
	if !ok {
		return nil, fmt.Errorf("%T output format is not supported", val)
	}
	return metrics, nil
}

func RangeDuration(args *Args, query string, baseTime time.Time, sine time.Duration, step time.Duration) (model.Value, error) {
	return Range(args, query, baseTime.Add(-sine), baseTime, step)
}

func ParseTime(timeStr string) (time.Time, error) {
	return time.Parse(time.RFC3339, timeStr)
}

func Range(args *Args, query string, start, end time.Time, step time.Duration) (model.Value, error) {
	res, err := queryBase.Range(&args.Args, query, start, end, step)
	if err != nil {
		return nil, fmt.Errorf("prometheus query range failed: %w", err)
	}

	return res, nil
}
