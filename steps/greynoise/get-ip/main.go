package main

import (
	"github.com/stackpulse/steps-sdk-go/env"
	"github.com/stackpulse/steps-sdk-go/step"
)

type GreynoiseGetIp struct {
    Ip string `env:"IP,required"`
GnApiKey string `env:"GN_API_KEY"`
}

type output struct {
    Name string `json:"name"`
Noise boolean `json:"noise"`
Riot boolean `json:"riot"`
Classification string `json:"classification"`
	step.Outputs
}

func (s *GreynoiseGetIp) Init() error {
	err := env.Parse(s)
	if err != nil {
		return err
	}

	return nil
}

func (s *GreynoiseGetIp) Run() (exitCode int, output []byte, err error) {
    // Replace This Comment with Step Logic

	return step.ExitCodeOK, []byte{}, nil
}

func main() {
	step.Run(&GreynoiseGetIp{})
}
