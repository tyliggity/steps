package main

import (
	"context"
	"encoding/json"

	"github.com/stackpulse/public-steps/jenkins/base"
	"github.com/stackpulse/steps-sdk-go/env"
	"github.com/stackpulse/steps-sdk-go/step"
)

//JenkinsStopBuild models input args for stopping an in-progress Jenkins build
type JenkinsStopBuild struct {
	Username string `env:"USERNAME,required"`
	Password string `env:"PASSWORD,required"`
	Host     string `env:"HOST,required"`
	JobName  string `env:"JOB_NAME,required"`
	BuildID  int    `env:"BUILD_ID,required"`
}

type output struct {
	Success bool `json:"success"`
	step.Outputs
}

//Init prepares step env
func (s *JenkinsStopBuild) Init() error {
	err := env.Parse(s)
	if err != nil {
		return err
	}

	return nil
}

//Run executes step logic
func (s *JenkinsStopBuild) Run() (int, []byte, error) {
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
	success, err := build.Stop(ctx)
	if err != nil {
		return step.ExitCodeFailure, nil, err
	}

	buildOutput := output{Success: success}

	outputJSON, err := json.Marshal(buildOutput)
	if err != nil {
		return step.ExitCodeFailure, nil, err
	}

	//success
	return step.ExitCodeOK, outputJSON, nil
}

func main() {
	step.Run(&JenkinsStopBuild{})
}
