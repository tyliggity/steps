package main

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/stackpulse/steps-sdk-go/step"
	"github.com/stackpulse/steps/steps/rundeck/base"
)

const (
	// https://docs.rundeck.com/docs/api/rundeck-api.html#running-a-job
	requestPath = "job/%v/executions"
	httpMethod  = http.MethodPost
)

type Args struct {
	base.Args
	JobID            string `env:"JOB_ID,required"`
	ExecutionOptions string `env:"OPTIONS" envDefault:"{}"`
}
type Outputs struct {
	step.Outputs
	JobID           string `json:"job_id"`
	JobName         string `json:"job_name"`
	ExecutionID     int    `json:"execution_id"`
	ExecutionStatus string `json:"execution_status"`
}

type ApiJsonResponse struct {
	base.ExecutionIDFromInt
}
type RundeckRunJob struct {
	*base.RundeckClient
	args Args
}

func (a Args) buildURI() string {
	return fmt.Sprintf(requestPath, a.JobID)
}

func (a Args) buildContent() ([]byte, error) {
	contentParams := fmt.Sprintf("{\"options\":%v}", a.ExecutionOptions)
	// TODO: add more parameters
	return []byte(contentParams), nil
}

func buildOutput(apiObject []byte) (Outputs, error) {
	var apiObjectStruct ApiJsonResponse
	err := json.Unmarshal(apiObject, &apiObjectStruct)
	if err != nil {
		return Outputs{}, fmt.Errorf("Failed to parse response. %v", err)
	}

	outputs := Outputs{
		Outputs:         step.Outputs{Object: apiObjectStruct},
		JobID:           apiObjectStruct.Job.ID,
		JobName:         apiObjectStruct.Job.Name,
		ExecutionID:     apiObjectStruct.ID,
		ExecutionStatus: apiObjectStruct.Status,
	}

	return outputs, nil
}

func (r *RundeckRunJob) Init() error {
	var args Args
	baseClient, err := base.NewRundeckClient(&args, httpMethod)
	if err != nil {
		return err
	}

	r.RundeckClient = baseClient
	r.args = args

	return nil
}

func (r *RundeckRunJob) Run() (int, []byte, error) {
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
	step.Run(&RundeckRunJob{})
}
