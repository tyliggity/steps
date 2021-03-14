package main

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/stackpulse/public-steps/common/step"
	"github.com/stackpulse/public-steps/steps/rundeck/base"
)

const (
	// https://docs.rundeck.com/docs/api/rundeck-api.html#aborting-executions
	requestPath = "execution/%v/abort"
	httpMethod  = http.MethodGet
)

type Args struct {
	base.Args
	ExecutionID string `env:"EXECUTION_ID,required"`
}
type Outputs struct {
	step.Outputs
	AbortionStatus  string `json:"abortion_status"`
	Reason          string `json:"reason,omitempty"`
	ExecutionID     int    `json:"execution_id"`
	ExecutionStatus string `json:"execution_status"`
}

type RundeckAbortExecution struct {
	*base.RundeckClient
	args Args
}

type ApiJsonResponse struct {
	Execution base.ExecutionIDFromString `json:"execution"`
	Abortion  base.Abort                 `json:"abort"`
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
		Outputs:         step.Outputs{Object: apiObjectStruct},
		AbortionStatus:  apiObjectStruct.Abortion.Status,
		Reason:          apiObjectStruct.Abortion.Reason,
		ExecutionID:     apiObjectStruct.Execution.ID,
		ExecutionStatus: apiObjectStruct.Execution.Status,
	}

	return outputs, nil
}

func (r *RundeckAbortExecution) Init() error {
	var args Args
	baseClient, err := base.NewRundeckClient(&args, httpMethod)
	if err != nil {
		return err
	}

	r.RundeckClient = baseClient
	r.args = args

	return nil
}

func (r *RundeckAbortExecution) Run() (int, []byte, error) {
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
	step.Run(&RundeckAbortExecution{})
}
