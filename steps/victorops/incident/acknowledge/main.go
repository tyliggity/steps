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

type IncidentAcknowledge struct {
	base.Args
	IncidentNumbers []string `env:"INCIDENT_NUMBERS,required"`
	UserName        string   `env:"USER_NAME,required"`
	AckMessage      string   `env:"ACK_MESSAGE"`
}

type ackRequest struct {
	IncidentNumbers []string `json:"incidentNames"`
	UserName        string   `json:"userName"`
	AckMessage      string   `json:"message"`
}

type ackResponse struct {
	Results []interface{} `json:"results"`
}

type output struct {
	IncidentNumbers []string `json:"incident_numbers"`
	EntityIds       []string `json:"entity_ids"`
	step.Outputs
}

func (s *IncidentAcknowledge) Init() error {
	err := env.Parse(s)
	if err != nil {
		return err
	}

	return nil
}

func (s *IncidentAcknowledge) Run() (int, []byte, error) {
	//prepare request body
	requestBodyData, err := json.Marshal(ackRequest{
		IncidentNumbers: s.IncidentNumbers,
		AckMessage:      s.AckMessage,
		UserName:        s.UserName,
	})
	if err != nil {
		return step.ExitCodeFailure, nil, err
	}
	requestBody := strings.NewReader(string(requestBodyData))

	//prepare request
	request, err := http.NewRequest(http.MethodPatch, "https://api.victorops.com/api-public/v1/incidents/ack", requestBody)
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
		return step.ExitCodeFailure, nil, fmt.Errorf("failed to acknowledge incident. got response code: %d: %s", resp.StatusCode, string(rawRespBody))
	}

	var responseBody ackResponse
	err = json.Unmarshal(rawRespBody, &responseBody)
	if err != nil {
		return step.ExitCodeFailure, nil, err
	}

	//prepare outputs
	incidentNumbers, entityIDs := []string{}, []string{}
	for _, res := range responseBody.Results {
		if resObj := res.(map[string]interface{}); resObj["cmdAccepted"].(bool) {
			incidentNumbers = append(incidentNumbers, resObj["incidentNumber"].(string))
			entityIDs = append(entityIDs, resObj["entityId"].(string))
		}
	}

	stepOutput := output{
		IncidentNumbers: incidentNumbers,
		EntityIds:       entityIDs,
		Outputs:         step.Outputs{Object: responseBody},
	}
	jsonOutput, err := json.Marshal(stepOutput)
	if err != nil {
		return step.ExitCodeFailure, nil, err
	}

	return step.ExitCodeOK, jsonOutput, nil
}

func main() {
	step.Run(&IncidentAcknowledge{})
}
