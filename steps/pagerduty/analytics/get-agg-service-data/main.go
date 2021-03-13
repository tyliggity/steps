package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/stackpulse/steps/common/env"
	"github.com/stackpulse/steps/common/log"
	"github.com/stackpulse/steps/common/step"
	"github.com/stackpulse/steps/pagerduty/base"
)

const (
	pagerdutyAnalyticsAPIUrl = "https://api.pagerduty.com/analytics/metrics/incidents/services"
)

type pagerDutyGetAggServiceData struct {
	base.Args
	TimeRange      string   `env:"TIME_RANGE,required"`
	ServiceIDs     []string `env:"SERVICE_IDS" envSeparator:","`
	ServiceNames   []string `env:"SERVICE_NAMES" envSeparator:","`
	CreatedAtStart string   `env:"CREATED_AT_START"`
	CreatedAtEnd   string   `env:"CREATED_AT_END"`
	AggregateUnit  string   `env:"AGGREGATE_UNIT"`
	TimeZone       string   `env:"TIME_ZONE"`
}

type output struct {
	step.Outputs
}

func (p *pagerDutyGetAggServiceData) Init() error {
	err := env.Parse(p)
	if err != nil {
		return err
	}

	return nil
}

func (p *pagerDutyGetAggServiceData) Run() (int, []byte, error) {
	content, err := p.buildContent()
	if err != nil {
		return step.ExitCodeFailure, nil, fmt.Errorf("build content: %w", err)
	}

	log.Debug("request content: %s", content)

	request, err := http.NewRequest(http.MethodPost, pagerdutyAnalyticsAPIUrl, bytes.NewReader(content))
	if err != nil {
		return step.ExitCodeFailure, nil, fmt.Errorf("create request: %w", err)
	}

	headers := map[string]string{
		"X-EARLY-ACCESS": "analytics-v2",
		"Accept":         "application/vnd.pagerduty+json;version=2",
		"Content-Type":   "application/json",
		"Authorization":  "Token token=" + p.PdToken,
	}
	for name, value := range headers {
		request.Header.Add(name, value)
	}

	resp, err := (&http.Client{}).Do(request)
	if err != nil {
		return step.ExitCodeFailure, nil, fmt.Errorf("send request: %w", err)
	}
	defer func() {
		if err := resp.Body.Close(); err != nil {
			log.Debug("Failed to close body: %v", err)
		}
	}()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return step.ExitCodeFailure, nil, fmt.Errorf("read all body: %w", err)
	}
	log.Debug("response received: '%s'", body)

	out := &output{}
	err = json.Unmarshal(body, &out.Object)
	if err != nil {
		return step.ExitCodeFailure, nil, fmt.Errorf("unmarshal body: %w", err)
	}

	jsonOutput, err := json.Marshal(out)
	if err != nil {
		return step.ExitCodeFailure, nil, fmt.Errorf("marshal output: %w", err)
	}

	return step.ExitCodeOK, jsonOutput, nil
}

func (p *pagerDutyGetAggServiceData) buildContent() ([]byte, error) {
	content := make(map[string]interface{})

	switch p.TimeRange {
	case "last_week":
		content["created_at_start"] = time.Now().AddDate(0, 0, -7).Format(time.RFC3339)
		content["created_at_end"] = time.Now().Format(time.RFC3339)
	case "last_month":
		content["created_at_start"] = time.Now().AddDate(0, -1, 0).Format(time.RFC3339)
		content["created_at_end"] = time.Now().Format(time.RFC3339)
	case "custom":
		content["created_at_start"] = p.CreatedAtStart
		content["created_at_end"] = p.CreatedAtEnd
	default:
		return nil, fmt.Errorf("unsupported time range: %s", p.TimeRange)
	}

	if len(p.ServiceNames) != 0 {
		idsByNames, err := base.ServiceIdsByNames(p.Args, p.ServiceNames)
		if err != nil {
			return nil, err
		}
		p.ServiceIDs = append(p.ServiceIDs, idsByNames...)
	}

	if len(p.ServiceIDs) != 0 {
		content["service_ids"] = p.ServiceIDs
	}

	filters := map[string]interface{}{"filters": content}

	if p.AggregateUnit != "" {
		filters["aggregate_unit"] = p.AggregateUnit
	}

	if p.TimeZone != "" {
		filters["time_zone"] = p.TimeZone
	}

	return json.Marshal(filters)
}

func main() {
	step.Run(&pagerDutyGetAggServiceData{})
}
