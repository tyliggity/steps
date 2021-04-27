package main

import (
	"github.com/stackpulse/steps-sdk-go/env"
	"github.com/stackpulse/steps-sdk-go/step"
)

type VirustotalGetIp struct {
    VtApiKey string `env:"VT_API_KEY,required"`
Ip string `env:"IP,required"`
}

type output struct {
    AsOwner string `json:"as_owner"`
Country string `json:"country"`
Reputation integer `json:"reputation"`
LastAnalysisStats json `json:"last_analysis_stats"`
	step.Outputs
}

func (s *VirustotalGetIp) Init() error {
	err := env.Parse(s)
	if err != nil {
		return err
	}

	return nil
}

func (s *VirustotalGetIp) Run() (exitCode int, output []byte, err error) {
    // Replace This Comment with Step Logic

	return step.ExitCodeOK, []byte{}, nil
}

func main() {
	step.Run(&VirustotalGetIp{})
}
