package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/stackpulse/public-steps/common/env"
	"github.com/stackpulse/public-steps/common/step"
	"github.com/stackpulse/public-steps/victorops/base"
)

type IncidentCreateNote struct {
	base.Args
	IncidentNumber int    `env:"INCIDENT_NUMBER,required"`
	NoteName       string `env:"NOTE_NAME,required"`
	NoteMessage    string `env:"NOTE_MESSAGE,required"`
}

type createRequest struct {
	NoteName    string               `json:"name"`
	NoteMessage createRequestMessage `json:"json_value"`
}

type createRequestMessage struct {
	Message string `json:"message"`
}

type output struct {
	Name string `json:"name"`
	step.Outputs
}

func (s *IncidentCreateNote) Init() error {
	err := env.Parse(s)
	if err != nil {
		return err
	}

	return nil
}

func (s *IncidentCreateNote) Run() (int, []byte, error) {
	//prepare request body
	requestBodyData, err := json.Marshal(createRequest{
		NoteName:    s.NoteName,
		NoteMessage: createRequestMessage{s.NoteMessage},
	})
	if err != nil {
		return step.ExitCodeFailure, nil, err
	}
	requestBody := strings.NewReader(string(requestBodyData))

	//prepare request
	request, err := http.NewRequest(
		http.MethodPost,
		fmt.Sprintf("https://api.victorops.com/api-public/v1/incidents/%d/notes", s.IncidentNumber),
		requestBody,
	)
	if err != nil {
		return step.ExitCodeFailure, nil, err
	}
	request.Header.Set("X-VO-Api-Id", s.ApiID)
	request.Header.Set("X-VO-Api-Key", s.ApiKey)
	request.Header.Set("Content-Type", "application/json")

	//send request
	client := http.Client{}
	resp, err := client.Do(request)
	if err != nil {
		return step.ExitCodeFailure, nil, err
	}

	//read response
	defer resp.Body.Close()
	rawRespBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return step.ExitCodeFailure, nil, err
	}

	//check HTTP error
	if resp.StatusCode != http.StatusOK {
		return step.ExitCodeFailure, nil, fmt.Errorf("failed to create note. got response code: %d: %s", resp.StatusCode, string(rawRespBody))
	}

	var responseBody map[string]interface{}
	err = json.Unmarshal(rawRespBody, &responseBody)
	if err != nil {
		return step.ExitCodeFailure, nil, err
	}

	//prepare output
	stepOutput := output{
		Name:    s.NoteName,
		Outputs: step.Outputs{Object: responseBody},
	}
	jsonOutput, err := json.Marshal(stepOutput)
	if err != nil {
		return step.ExitCodeFailure, nil, err
	}

	return step.ExitCodeOK, jsonOutput, nil
}

func main() {
	step.Run(&IncidentCreateNote{})
}
