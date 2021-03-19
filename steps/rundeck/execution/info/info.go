package main

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/stackpulse/public-steps/steps/rundeck/base"
	"github.com/stackpulse/steps-sdk-go/step"
)

const (
	// https://docs.rundeck.com/docs/api/rundeck-api.html#execution-info
	requestPath = "execution/%v"
	httpMethod  = http.MethodGet
)

type Args struct {
	base.Args
	ExecutionID string `env:"EXECUTION_ID,required"`
}
type Outputs struct {
	base.ExecutionIDFromInt
}

type RundeckExecutionInfo struct {
	*base.RundeckClient
	args Args
}

type ApiJsonResponse struct {
	base.ExecutionIDFromInt
}

func (a Args) buildURI() string {
	return fmt.Sprintf(requestPath, a.ExecutionID)
}

func buildOutput(apiObject []byte) (Outputs, error) {
	var apiObjectStruct ApiJsonResponse
	err := json.Unmarshal(apiObject, &apiObjectStruct)
	if err != nil {
		return Outputs{}, fmt.Errorf("Failed to parse response. %v", err)
	}

	outputs := Outputs{
		ExecutionIDFromInt: apiObjectStruct.ExecutionIDFromInt,
	}

	return outputs, nil
}

func (r *RundeckExecutionInfo) Init() error {
	var args Args
	baseClient, err := base.NewRundeckClient(&args, httpMethod)
	if err != nil {
		return err
	}

	r.RundeckClient = baseClient
	r.args = args

	return nil
}

func (r *RundeckExecutionInfo) Run() (int, []byte, error) {
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
	step.Run(&RundeckExecutionInfo{})
}
