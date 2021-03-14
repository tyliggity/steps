package main

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/stackpulse/public-steps/common/log"
	"github.com/stackpulse/public-steps/common/step"
	"github.com/stackpulse/public-steps/steps/jenkins/base"
)

const MaxJobWait = 185 * time.Second // Max timeout for step is 200 seconds, so I want to return earlier
const SleepBetweenBuildCheck = 3 * time.Second
const ApiCallTimeout = 10 * time.Second

type Args struct {
	base.Args
	JobName    string `env:"JOB_NAME,required"`
	Params     string `env:"PARAMS"`
	ShouldWait bool   `env:"WAIT_FOR_BUILD" envDefault:"false"`
}

type Output struct {
	BuildQueueID int64  `json:"build_queue_id"`
	BuildID      string `json:"build_id"`
	BuildOutput  string `json:"build_output,omitempty"`
	Succeeded    bool   `json:"succeeded,omitempty"`
}

type JenkinsRun struct {
	*base.JenkinsBaseCommand
	args Args
}

func timeoutContext() context.Context {
	ctx, _ := context.WithTimeout(context.Background(), ApiCallTimeout)
	return ctx
}

func (j *JenkinsRun) Init() error {
	baseCommand, err := base.NewBaseJenkinsCommand(&j.args)
	if err != nil {
		return err
	}

	j.JenkinsBaseCommand = baseCommand
	return nil
}

func (j *JenkinsRun) marshalAndReturn(output *Output, forceErr error) (int, []byte, error) {
	marshaledOutput, err := json.Marshal(output)
	if err != nil {
		return step.ExitCodeFailure, nil, fmt.Errorf("marshal output: %w", err)
	}

	var retErr error
	exitCode := step.ExitCodeOK
	if forceErr != nil {
		retErr = forceErr
		exitCode = step.ExitCodeFailure
	}
	return exitCode, marshaledOutput, retErr
}

func (j *JenkinsRun) Run() (int, []byte, error) {
	var params map[string]string
	if j.args.Params != "" {
		if err := json.Unmarshal([]byte(j.args.Params), &params); err != nil {
			return step.ExitCodeFailure, nil, fmt.Errorf("PARAMS must be a valid {key:value} json")
		}
	}

	log.Debug("Creating build job for build %s", j.args.JobName)
	buildQueueID, err := j.Client().BuildJob(timeoutContext(), j.args.JobName, params)
	if err != nil {
		return step.ExitCodeFailure, nil, fmt.Errorf("build jenkins job: %w", err)
	}

	buildID, err := j.Client().GetBuildFromQueueID(timeoutContext(), buildQueueID)
	if err != nil {
		return step.ExitCodeFailure, nil, fmt.Errorf("get build from queue id: %w", err)
	}

	output := &Output{
		BuildQueueID: buildQueueID,
		BuildID:      buildID.Raw.ID,
	}

	if !j.args.ShouldWait {
		return j.marshalAndReturn(output, nil)
	}

	build, err := j.Client().GetBuildFromQueueID(timeoutContext(), buildQueueID)
	if err != nil {
		return step.ExitCodeFailure, nil, fmt.Errorf("get build %d: %w", buildQueueID, err)
	}

	maxWaitTime := time.Now().Add(MaxJobWait)
	for build.IsRunning(timeoutContext()) {
		if time.Now().After(maxWaitTime) {
			return j.marshalAndReturn(output, fmt.Errorf("timeout waiting for build to finished"))
		}
		log.Debugln("Build is still running..")
		time.Sleep(SleepBetweenBuildCheck)
		if _, err := build.Poll(timeoutContext()); err != nil {
			return j.marshalAndReturn(output, fmt.Errorf("poll for build %d: %w", buildQueueID, err))
		}
	}

	log.Debug("Build finished")
	output.BuildOutput = build.GetConsoleOutput(timeoutContext())
	output.Succeeded = build.IsGood(timeoutContext())
	return j.marshalAndReturn(output, nil)
}

func main() {
	step.Run(&JenkinsRun{})
}
