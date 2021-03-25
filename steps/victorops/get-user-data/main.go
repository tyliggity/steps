package main

import (
	"encoding/json"

	envconf "github.com/caarlos0/env/v6"
	"github.com/stackpulse/steps-sdk-go/step"
	"github.com/stackpulse/steps/victorops/base"
	"github.com/victorops/go-victorops/victorops"
)

type Args struct {
	base.Args
	UserName string `env:"USER_NAME,required"`
}

type VictorOpsGetUserData struct {
	args Args
}

func (v *VictorOpsGetUserData) Init() error {
	return envconf.Parse(&v.args)
}

func (v *VictorOpsGetUserData) Run() (exitCode int, output []byte, err error) {
	victoropsClient := victorops.NewClient(v.args.ApiID, v.args.ApiKey, "https://api.victorops.com")

	user, _, err := victoropsClient.GetUser(v.args.UserName)
	if err != nil {
		return step.ExitCodeFailure, nil, err
	}

	marshaledUser, err := json.Marshal(user)
	if err != nil {
		return step.ExitCodeFailure, nil, err
	}

	return step.ExitCodeOK, marshaledUser, nil
}

func main() {
	step.Run(&VictorOpsGetUserData{})
}
