package main

import (
	"io/ioutil"
	"net/http"

	envconf "github.com/caarlos0/env/v6"
	"github.com/stackpulse/steps-sdk-go/step"
	"github.com/stackpulse/steps/victorops/base"
)

type Args struct {
	base.Args
}

type VictorOpsGetOncall struct {
	args Args
}

func (v *VictorOpsGetOncall) Init() error {
	return envconf.Parse(&v.args)
}

func (v *VictorOpsGetOncall) Run() (exitCode int, output []byte, err error) {
	client := http.Client{}
	request, err := http.NewRequest(http.MethodGet, "https://api.victorops.com/api-public/v1/oncall/current", nil)
	if err != nil {
		return step.ExitCodeFailure, nil, err
	}

	request.Header.Set("X-VO-Api-Id", v.args.ApiID)
	request.Header.Set("X-VO-Api-Key", v.args.ApiKey)
	request.Header.Set("Content-Type", "application/json")

	resp, err := client.Do(request)
	if err != nil {
		return step.ExitCodeFailure, nil, err
	}

	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return step.ExitCodeFailure, nil, err
	}

	return step.ExitCodeOK, body, nil
}

func main() {
	step.Run(&VictorOpsGetOncall{})
}
