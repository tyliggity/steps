package main

import (
	"net/http"

	"github.com/Jeffail/gabs/v2"
	envconf "github.com/caarlos0/env/v6"
	"github.com/opsgenie/opsgenie-go-sdk-v2/client"
	"github.com/opsgenie/opsgenie-go-sdk-v2/incident"
	"github.com/pkg/errors"

	"github.com/stackpulse/public-steps/common/step"
	"github.com/stackpulse/public-steps/opsgenie/base"
)

type Args struct {
	base.Args
	Id   string `env:"INCIDENT_ID,required"`
	Note string `env:"NOTE,required"`
}

type OpsgenieIncidentResolve struct {
	args Args
}

func (o *OpsgenieIncidentResolve) Init() error {
	err := envconf.Parse(&o.args)
	if err != nil {
		return err
	}

	return nil
}

func (o *OpsgenieIncidentResolve) Run() (exitCode int, output []byte, err error) {
	gc := gabs.New()
	ogClient, err := client.NewOpsGenieClient(base.Config(o.args.Args))
	if err != nil {
		return step.ExitCodeFailure, nil, err
	}

	result := &incident.AsyncResult{}
	err = ogClient.Exec(nil, &ResolveRequest{
		Id:   o.args.Id,
		Note: o.args.Note,
	}, result)
	if err != nil {
		return step.ExitCodeFailure, nil, err
	}

	gc.Set(true, "success")

	return step.ExitCodeOK, gc.Bytes(), nil
}

func main() {
	step.Run(&OpsgenieIncidentResolve{})
}

type ResolveRequest struct {
	client.BaseRequest
	Id         string
	Identifier incident.IdentifierType
	Note       string `json:"note,omitempty"`
}

func (r *ResolveRequest) Validate() error {
	if r.Id == "" {
		return errors.New("Incident ID cannot be blank.")
	}
	if r.Identifier != "" && r.Identifier != incident.Id && r.Identifier != incident.Tiny {
		return errors.New("Identifier type should be one of these: 'Id', 'Tiny' or empty.")
	}
	return nil
}

func (r *ResolveRequest) ResourcePath() string {
	return "/v1/incidents/" + r.Id + "/resolve"
}

func (r *ResolveRequest) Method() string {
	return http.MethodPost
}

func (r *ResolveRequest) RequestParams() map[string]string {
	params := make(map[string]string)

	if r.Identifier == incident.Tiny {
		params["identifierType"] = "tiny"
	} else {
		params["identifierType"] = "id"
	}

	return params
}
