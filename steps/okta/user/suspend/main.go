package main

import (
	"github.com/stackpulse/steps-sdk-go/env"
	"github.com/stackpulse/steps-sdk-go/step"
)

type UserSuspend struct {
    OktaApiToken string `env:"OKTA_API_TOKEN,required"`
OktaDomain string `env:"OKTA_DOMAIN,required"`
UserId string `env:"USER_ID,required"`
}

type output struct {
    Result string `json:"result"`
	step.Outputs
}

func (s *UserSuspend) Init() error {
	err := env.Parse(s)
	if err != nil {
		return err
	}

	return nil
}

func (s *UserSuspend) Run() (exitCode int, output []byte, err error) {
    // Replace This Comment with Step Logic

	return step.ExitCodeOK, []byte{}, nil
}

func main() {
	step.Run(&UserSuspend{})
}
