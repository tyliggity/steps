package main

import (
	"github.com/stackpulse/steps-sdk-go/env"
	"github.com/stackpulse/steps-sdk-go/step"
)

type VirustotalGetUrl struct {
    VtApiKey string `env:"VT_API_KEY,required"`
Url string `env:"URL,required"`
}

type output struct {
    Stats json `json:"stats"`
	step.Outputs
}

func (s *VirustotalGetUrl) Init() error {
	err := env.Parse(s)
	if err != nil {
		return err
	}

	return nil
}

func (s *VirustotalGetUrl) Run() (exitCode int, output []byte, err error) {
    // Replace This Comment with Step Logic

	return step.ExitCodeOK, []byte{}, nil
}

func main() {
	step.Run(&VirustotalGetUrl{})
}
