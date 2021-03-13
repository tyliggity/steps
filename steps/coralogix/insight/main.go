package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	envconf "github.com/caarlos0/env/v6"
	"github.com/stackpulse/public-steps/common/log"
	"github.com/stackpulse/public-steps/common/step"
	"github.com/stackpulse/public-steps/coralogix/base"
)

type Args struct {
	base.Args
	Type      string `env:"TYPE,required"`
	AppName   string `env:"APPLICATION_NAME"`
	SubName   string `env:"SUBSYSTEM_NAME"`
	TagName   string `env:"TAG_NAME"`
	Severity  string `env:"SEVERITY"`
	StartDate int64  `env:"START_DATE" envDefault:"0"`
	EndDate   int64  `env:"END_DATE" envDefault:"0"`
}

const (
	coralogixInsightsAPIUrl = "https://api.coralogix.com/api/v1/external/insights"
	coralogixTokenHeader    = "Authorization"
	coralogixBearerTokenFmt = "Bearer %s"
)

type CoralogixInsight struct {
	args Args
}

func (l *CoralogixInsight) Init() error {
	err := envconf.Parse(&l.args)
	if err != nil {
		return err
	}

	return nil
}

func (l *CoralogixInsight) buildContent() ([]byte, error) {
	var content map[string]interface{}

	queryJson := fmt.Sprintf("{\"type\": \"%s\"}", l.args.Type)
	err := json.Unmarshal([]byte(queryJson), &content)
	if err != nil {
		return nil, fmt.Errorf("failed to build body: %w", err)
	}

	if l.args.TagName != "" {
		if l.args.StartDate != 0 || l.args.EndDate != 0 {
			return nil, fmt.Errorf("START_DATE or END_DATE cannot be set with TAG_NAME")
		}
		content["tagName"] = l.args.TagName
	} else {
		if l.args.StartDate == 0 || l.args.EndDate == 0 {
			return nil, fmt.Errorf("START_DATE or END_DATE must be set if TAG_NAME is undefined")
		}

		content["startDate"] = l.args.StartDate
		content["endDate"] = l.args.EndDate
	}

	if l.args.Severity != "" {
		content["severity"] = l.args.Severity
	}

	if l.args.AppName != "" {
		content["applicationName"] = l.args.AppName
	}

	if l.args.SubName != "" {
		content["subsystemName"] = l.args.SubName
	}

	return json.Marshal(content)
}

func (l *CoralogixInsight) Run() (int, []byte, error) {
	client := &http.Client{}

	content, err := l.buildContent()
	if err != nil {
		return 1, nil, err
	}

	log.Debug("request content: %s", string(content))

	request, err := http.NewRequest(http.MethodPost, coralogixInsightsAPIUrl, bytes.NewReader(content))
	if err != nil {
		return 1, nil, err
	}

	request.Header.Add(coralogixTokenHeader, fmt.Sprintf(coralogixBearerTokenFmt, l.args.Token))
	request.Header.Add("Content-Type", "application/json")

	resp, err := client.Do(request)
	if err != nil {
		return 1, nil, err
	}

	defer resp.Body.Close()
	output, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return 1, nil, err
	}

	return 0, output, nil
}

func main() {
	step.Run(&CoralogixInsight{})
}
