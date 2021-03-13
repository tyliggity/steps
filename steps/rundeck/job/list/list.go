package main

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/stackpulse/public-steps/common/step"
	"github.com/stackpulse/public-steps/public-steps/rundeck/base"
)

const (
	// https://docs.rundeck.com/docs/api/rundeck-api.html#listing-jobs
	requestPath = "project/%v/jobs"
	httpMethod  = http.MethodGet
)

type Args struct {
	base.Args
	ProjectName string `env:"PROJECT,required"`
}

type Outputs struct {
	Jobs []base.Job `json:"jobs"`
}

type ApiJsonResponse []base.Job

type RundeckListJobs struct {
	*base.RundeckClient
	args Args
}

func (a Args) buildURI() string {
	return fmt.Sprintf(requestPath, a.ProjectName)
}

func (a Args) buildContent() ([]byte, error) {
	// TODO: add more parameters
	return []byte{}, nil
}

func buildOutput(apiObject []byte) (Outputs, error) {
	var apiObjectStruct ApiJsonResponse
	err := json.Unmarshal(apiObject, &apiObjectStruct)
	if err != nil {
		return Outputs{}, fmt.Errorf("Failed to parse response. %v", err)
	}

	outputs := Outputs{
		Jobs: apiObjectStruct,
	}

	return outputs, nil
}

func (r *RundeckListJobs) Init() error {
	var args Args
	baseClient, err := base.NewRundeckClient(&args, httpMethod)
	if err != nil {
		return err
	}

	r.RundeckClient = baseClient
	r.args = args

	return nil
}

func (r *RundeckListJobs) Run() (int, []byte, error) {
	uri := r.args.buildURI()
	contentParams, err := r.args.buildContent()
	if err != nil {
		return step.ExitCodeFailure, nil, err
	}

	errCode, resBody, err := r.MakeRequest(uri, httpMethod, contentParams)
	if err != nil {
		return errCode, resBody, err
	}

	output, err := buildOutput(resBody)
	if err != nil {
		return step.ExitCodeFailure, resBody, err
	}

	byteOutput, err := json.Marshal(output)
	if err != nil {
		return step.ExitCodeFailure, resBody, err
	}

	return step.ExitCodeOK, byteOutput, nil
}

func main() {
	step.Run(&RundeckListJobs{})
}
