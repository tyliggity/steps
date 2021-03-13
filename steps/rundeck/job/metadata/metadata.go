package main

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/stackpulse/public-steps/common/step"
	"github.com/stackpulse/public-steps/public-steps/rundeck/base"
)

const (
	// https://docs.rundeck.com/docs/api/rundeck-api.html#get-job-metadata
	requestPath = "job/%v/info"
	httpMethod  = http.MethodGet
)

type Args struct {
	base.Args
	JobID string `env:"JOB_ID,required"`
}
type Outputs struct {
	base.Job
}

type ApiJsonResponse struct {
	base.Job
}
type RundeckJobMetadata struct {
	*base.RundeckClient
	args Args
}

func (a Args) buildURI() string {
	return fmt.Sprintf(requestPath, a.JobID)
}

func buildOutput(apiObject []byte) (Outputs, error) {
	var apiObjectStruct ApiJsonResponse
	err := json.Unmarshal(apiObject, &apiObjectStruct)
	if err != nil {
		return Outputs{}, fmt.Errorf("Failed to parse response. %v", err)
	}

	outputs := Outputs{
		Job: apiObjectStruct.Job,
	}

	return outputs, nil
}

func (r *RundeckJobMetadata) Init() error {
	var args Args
	baseClient, err := base.NewRundeckClient(&args, httpMethod)
	if err != nil {
		return err
	}

	r.RundeckClient = baseClient
	r.args = args

	return nil
}

func (r *RundeckJobMetadata) Run() (int, []byte, error) {
	uri := r.args.buildURI()
	errCode, resBody, err := r.MakeRequest(uri, httpMethod, []byte{})
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
	step.Run(&RundeckJobMetadata{})
}
