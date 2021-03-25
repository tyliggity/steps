package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	envconf "github.com/caarlos0/env/v6"
	"github.com/stackpulse/steps-sdk-go/log"
	"github.com/stackpulse/steps-sdk-go/step"
	"github.com/stackpulse/steps/logzio/base"
)

type Args struct {
	base.Args
	Query   string `env:"QUERY,required"`
	DaysOff string `env:"DAYS_OFF"`
	From    string `env:"FROM"`
	Size    string `env:"SIZE"`
}

const (
	logzioQueryUrlFormat = "%s/v1/search"
	logzioTokenHeader    = "X-API-TOKEN"
)

type LogzIoQuery struct {
	args Args
}

func (l *LogzIoQuery) Init() error {
	err := envconf.Parse(&l.args)
	if err != nil {
		return err
	}

	return nil
}

func (l *LogzIoQuery) buildUrl() string {
	url := logzioQueryUrlFormat
	if l.args.DaysOff != "" {
		url += fmt.Sprintf("?daysOff=%s", l.args.DaysOff)
	}

	return fmt.Sprintf(url, l.args.Endpoint)
}

func (l *LogzIoQuery) buildContent() ([]byte, error) {
	var content map[string]interface{}

	queryJson := fmt.Sprintf("{\"query\": %s}", l.args.Query)
	err := json.Unmarshal([]byte(queryJson), &content)
	if err != nil {
		return nil, err
	}

	if l.args.Size != "" {
		content["size"] = l.args.Size
	}

	if l.args.From != "" {
		content["from"] = l.args.From
	}

	return json.Marshal(content)
}

func (l *LogzIoQuery) Run() (int, []byte, error) {
	client := &http.Client{}

	content, err := l.buildContent()
	if err != nil {
		return 1, nil, err
	}

	log.Debug("request content: %s", string(content))

	request, err := http.NewRequest(http.MethodPost, l.buildUrl(), bytes.NewReader(content))
	if err != nil {
		return 1, nil, err
	}

	request.Header.Add(logzioTokenHeader, l.args.Token)
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
	step.Run(&LogzIoQuery{})
}
