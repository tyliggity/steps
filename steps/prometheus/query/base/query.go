package query

import (
	"context"
	"time"

	promv1 "github.com/prometheus/client_golang/api/prometheus/v1"
	"github.com/prometheus/common/model"
	"github.com/stackpulse/steps/prometheus/base"
)

type Args struct {
	base.Args
}

func Query(args *Args, query string, ts time.Time) (model.Value, error) {
	api, err := base.API(&args.Args)
	if err != nil {
		return nil, err
	}
	val, _, err := api.Query(context.Background(), query, ts)
	if err != nil {
		return nil, err
	}
	return val, err
}

func Range(args *Args, query string, start, end time.Time, step time.Duration) (model.Value, error) {
	api, err := base.API(&args.Args)
	if err != nil {
		return nil, err
	}
	val, _, err := api.QueryRange(context.Background(), query, promv1.Range{
		Start: start,
		End:   end,
		Step:  step,
	})

	if err != nil {
		return nil, err
	}
	return val, err
}
