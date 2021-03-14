package base

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"time"

	"github.com/stackpulse/public-steps/common/env"
	"github.com/stackpulse/public-steps/common/step"
)

const (
	rundeckBasePathFormat    = "%v://%v:%v/api/38/%v"
	rundeckAuthTokenHeader   = "X-Rundeck-Auth-Token"
	httpTimeout              = 10 * time.Second
	AcceptedFormatHeaderJSON = "application/json"
)

type Args struct {
	Host   string `env:"HOST,required"`
	Port   int    `env:"PORT,required"`
	Scheme string `env:"HTTP_SCHEME" envDefault:"https"`
	Token  string `env:"AUTH_TOKEN,required"`
}

func (a Args) BaseArgs() Args {
	return a
}

type BaseArgs interface {
	BaseArgs() Args
}

type RundeckClient struct {
	Args
	client *http.Client
}

func NewRundeckClient(args BaseArgs, httpMethod string) (*RundeckClient, error) {
	if err := env.Parse(args); err != nil {
		return nil, fmt.Errorf("parse args: %w", err)
	}

	httpClient := &http.Client{
		Timeout: httpTimeout,
	}

	return &RundeckClient{Args: args.BaseArgs(), client: httpClient}, nil
}

func (r RundeckClient) newHTTPRequest(httpMethod, uri string, contentParams []byte) (*http.Request, error) {
	urlString := fmt.Sprintf(rundeckBasePathFormat, r.Args.Scheme, r.Args.Host, r.Args.Port, uri)
	parsedUrl, err := url.Parse(urlString)
	if err != nil {
		return nil, fmt.Errorf("Invalid host, port, scheme or entity id format. url: %v", urlString)
	}

	req, err := http.NewRequest(httpMethod, parsedUrl.String(), bytes.NewBuffer(contentParams))
	if err != nil {
		return req, err
	}

	req.Header.Set("Content-Type", AcceptedFormatHeaderJSON)
	req.Header.Set("Accept", AcceptedFormatHeaderJSON)
	req.Header.Set(rundeckAuthTokenHeader, r.Args.Token)

	return req, nil
}

func extractResponse(res *http.Response) (int, []byte, error) {
	readBody, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return step.ExitCodeFailure, readBody, err
	}

	if res.StatusCode/100 != 2 {
		var jsonBody map[string]interface{}
		err := json.Unmarshal(readBody, &jsonBody)
		if msg, msgExists := jsonBody["message"].(string); err == nil && msgExists {
			return step.ExitCodeFailure, nil, fmt.Errorf(msg)
		}

		return step.ExitCodeFailure, []byte(res.Request.URL.String()), fmt.Errorf("A non-success response code was received from rundeck [%s]", res.Status)
	}

	serializedBody, err := jsonKeysCaseSerializer(readBody)
	if err != nil {
		return step.ExitCodeFailure, readBody, err
	}

	return step.ExitCodeOK, serializedBody, nil
}

func (r *RundeckClient) MakeRequest(uri, httpMethod string, contentParams []byte) (int, []byte, error) {
	httpRequest, err := r.newHTTPRequest(httpMethod, uri, contentParams)
	if err != nil {
		return step.ExitCodeFailure, nil, err
	}

	res, err := r.client.Do(httpRequest)
	if err != nil {
		return step.ExitCodeFailure, nil, err
	}
	defer res.Body.Close()

	return extractResponse(res)
}
