package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	envconf "github.com/caarlos0/env/v6"
	"github.com/stackpulse/steps/common/log"
	"github.com/stackpulse/steps/common/step"
	"github.com/stackpulse/steps/coralogix/base"
)

type Args struct {
	base.Args
	Query string `env:"QUERY,required"`
	From  string `env:"FROM"`
	Size  string `env:"SIZE"`
	Sort  string `env:"SORT"`
}

const (
	coralogixElasticsearchAPIUrl = "https://coralogix-esapi.coralogix.com:9443/*/_search"
	coralogixTokenHeader         = "token"
)

type CoralogixQuery struct {
	args Args
}

func (l *CoralogixQuery) Init() error {
	err := envconf.Parse(&l.args)
	if err != nil {
		return err
	}

	return nil
}

func (l *CoralogixQuery) buildContent() ([]byte, error) {
	var content map[string]interface{}

	queryJson := fmt.Sprintf("{\"query\": %s}", l.args.Query)
	err := json.Unmarshal([]byte(queryJson), &content)
	if err != nil {
		return nil, fmt.Errorf("failed to build body: %w", err)
	}

	if l.args.Size != "" {
		content["size"] = l.args.Size
	}

	if l.args.From != "" {
		content["from"] = l.args.From
	}

	if l.args.Sort != "" {
		content["sort"] = l.args.Sort
	}

	return json.Marshal(content)
}

func (l *CoralogixQuery) Run() (int, []byte, error) {
	client := &http.Client{}

	content, err := l.buildContent()
	if err != nil {
		return 1, nil, err
	}

	log.Debug("request content: %s", string(content))

	request, err := http.NewRequest(http.MethodPost, coralogixElasticsearchAPIUrl, bytes.NewReader(content))
	if err != nil {
		return 1, nil, err
	}

	request.Header.Add(coralogixTokenHeader, l.args.Token)
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
	step.Run(&CoralogixQuery{})
}
