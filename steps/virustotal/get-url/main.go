package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/stackpulse/steps-sdk-go/env"
	"github.com/stackpulse/steps-sdk-go/step"
)

type VirustotalGetUrl struct {
	VtApiKey string `env:"VT_API_KEY,required"`
	Url      string `env:"URL,required"`
}

type stepOutput struct {
	Stats statsOutput `json:"stats"`
	step.Outputs
}

type statsOutput struct {
	Harmless   int `json:"harmless"`
	Malicious  int `json:"malicious"`
	Suspicious int `json:"suspicious"`
	Timeout    int `json:"timeout"`
	Undetected int `json:"undetected"`
}

type apiResponse struct {
	Data struct {
		Attributes struct {
			LastAnalysisStats statsOutput `json:"last_analysis_stats"`
			AsOwner           string      `json:"as_owner"`
			Country           string      `json:"country"`
			Reputation        int         `json:"reputation"`
		} `json:"attributes"`
	} `json:"data"`
}

func (s *VirustotalGetUrl) Init() error {
	err := env.Parse(s)
	if err != nil {
		return err
	}

	return nil
}

func (s *VirustotalGetUrl) Run() (exitCode int, output []byte, err error) {
	//prepare request
	reqURL := fmt.Sprintf("https://www.virustotal.com/api/v3/domains/%s", s.Url)
	stepReq, err := http.NewRequest(http.MethodGet, reqURL, nil)
	if err != nil {
		return step.ExitCodeFailure, nil, err
	}
	stepReq.Header.Set("x-apikey", s.VtApiKey)

	//send request
	res, err := http.DefaultClient.Do(stepReq)
	if err != nil || res.StatusCode != http.StatusOK {
		return step.ExitCodeFailure, nil, err
	}

	//parse response
	defer res.Body.Close()
	rawRes, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return step.ExitCodeFailure, nil, err
	}

	var resp apiResponse
	err = json.Unmarshal(rawRes, &resp)
	if err != nil {
		return step.ExitCodeFailure, nil, err
	}

	//unmarshal full response
	var fullResp map[string]interface{}
	err = json.Unmarshal(rawRes, &fullResp)
	if err != nil {
		return step.ExitCodeFailure, nil, err
	}

	//prepare output
	out := stepOutput{
		Stats:   resp.Data.Attributes.LastAnalysisStats,
		Outputs: step.Outputs{Object: fullResp},
	}
	outputBytes, err := json.Marshal(out)
	if err != nil {
		return step.ExitCodeFailure, nil, err
	}

	//success
	return step.ExitCodeOK, outputBytes, nil
}

func main() {
	step.Run(&VirustotalGetUrl{})
}
