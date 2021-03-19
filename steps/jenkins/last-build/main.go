package main

import (
	"context"
	"encoding/json"

	"github.com/stackpulse/public-steps/jenkins/base"
	"github.com/stackpulse/steps-sdk-go/env"
	"github.com/stackpulse/steps-sdk-go/step"
)

//JenkinsLastBuild models input args for retrieving latest Jenkins build info
type JenkinsLastBuild struct {
	Username string `env:"USERNAME,required"`
	Password string `env:"PASSWORD,required"`
	Host     string `env:"HOST,required"`
	JobName  string `env:"JOB_NAME,required"`
}

type output struct {
	LastBuild           int `json:"last_build"`
	LastCompletedBuild  int `json:"last_completed_build"`
	LastFailedBuild     int `json:"last_failed_build"`
	LastSuccessfulBuild int `json:"last_successful_build"`
	LastStableBuild     int `json:"last_stable_build"`
	step.Outputs
}

//Init prepares step env
func (s *JenkinsLastBuild) Init() error {
	err := env.Parse(s)
	if err != nil {
		return err
	}

	return nil
}

//Run executes step logic
func (s *JenkinsLastBuild) Run() (int, []byte, error) {
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
	job, err := jenkins.GetJob(ctx, s.JobName)
	if err != nil {
		return step.ExitCodeFailure, nil, err
	}

	lastBuild, err := job.GetLastBuild(ctx)
	if err != nil {
		return step.ExitCodeFailure, nil, err
	}
	lastCompletedBuild, err := job.GetLastCompletedBuild(ctx)
	if err != nil {
		return step.ExitCodeFailure, nil, err
	}
	lastFailedBuild, err := job.GetLastFailedBuild(ctx)
	if err != nil {
		return step.ExitCodeFailure, nil, err
	}
	lastSuccessfulBuild, err := job.GetLastSuccessfulBuild(ctx)
	if err != nil {
		return step.ExitCodeFailure, nil, err
	}
	lastStableBuild, err := job.GetLastStableBuild(ctx)
	if err != nil {
		return step.ExitCodeFailure, nil, err
	}

	buildOutput := output{
		LastBuild:           int(lastBuild.GetBuildNumber()),
		LastCompletedBuild:  int(lastCompletedBuild.GetBuildNumber()),
		LastFailedBuild:     int(lastFailedBuild.GetBuildNumber()),
		LastSuccessfulBuild: int(lastSuccessfulBuild.GetBuildNumber()),
		LastStableBuild:     int(lastStableBuild.GetBuildNumber()),
		Outputs:             step.Outputs{Object: job.Raw},
	}

	outputJSON, err := json.Marshal(buildOutput)
	if err != nil {
		return step.ExitCodeFailure, nil, err
	}

	//success
	return step.ExitCodeOK, outputJSON, nil
}

func main() {
	step.Run(&JenkinsLastBuild{})
}
