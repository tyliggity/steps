package main

import (
	"context"
	"encoding/json"
	"time"

	"github.com/stackpulse/steps-sdk-go/env"
	"github.com/stackpulse/steps-sdk-go/step"
	"github.com/stackpulse/steps/jenkins/base"
)

//JenkinsGetBuild models input args for retrieving info about a specific Jenkins build
type JenkinsGetBuild struct {
	Username string `env:"USERNAME,required"`
	Password string `env:"PASSWORD,required"`
	Host     string `env:"HOST,required"`
	JobName  string `env:"JOB_NAME,required"`
	BuildID  int    `env:"BUILD_ID,required"`
}

type output struct {
	Result    string `json:"result"`
	Duration  int    `json:"duration"`
	Name      string `json:"name"`
	Timestamp string `json:"timestamp"`
	step.Outputs
}

//Init prepares step env
func (s *JenkinsGetBuild) Init() error {
	err := env.Parse(s)
	if err != nil {
		return err
	}

	return nil
}

//Run executes step logic
func (s *JenkinsGetBuild) Run() (int, []byte, error) {
	baseCommand, err := base.NewBaseJenkinsCommand(&base.Args{
		Username: s.Username,
		Password: s.Password,
		Host:     s.Host,
	})
	if err != nil {
		return step.ExitCodeFailure, nil, err
	}

	jenkins := baseCommand.Client()
	ctx := context.Background()
	build, err := jenkins.GetBuild(ctx, s.JobName, int64(s.BuildID))
	if err != nil {
		return step.ExitCodeFailure, nil, err
	}

	buildOutput := output{
		Result:    build.GetResult(),
		Duration:  int(build.GetDuration()),
		Name:      build.Raw.DisplayName,
		Timestamp: build.GetTimestamp().Format(time.RFC3339),
		Outputs:   step.Outputs{Object: build.Raw},
	}

	outputJSON, err := json.Marshal(buildOutput)
	if err != nil {
		return step.ExitCodeFailure, nil, err
	}

	//success
	return step.ExitCodeOK, outputJSON, nil
}

func main() {
	step.Run(&JenkinsGetBuild{})
}
