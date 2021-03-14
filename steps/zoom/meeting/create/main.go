package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	envconf "github.com/caarlos0/env/v6"
	"github.com/himalayan-institute/zoom-lib-golang"
	"github.com/stackpulse/public-steps/common/step"
	"github.com/stackpulse/public-steps/zoom/base"
	"golang.org/x/oauth2"
	"io/ioutil"
	"net/http"
)

type Args struct {
	base.Args
	HostEmail     string `env:"HOST_EMAIL,required"`
	MeetingTopic  string `env:"MEETING_TOPIC"`
	ZoomApiKey    string `env:"ZOOM_API_KEY"`
	ZoomApiSecret string `env:"ZOOM_API_SECRET"`
}

type ZoomMeetingCreate struct {
	args  Args
	token *oauth2.Token
}

func (z *ZoomMeetingCreate) Init() error {
	if err := envconf.Parse(&z.args); err != nil {
		return fmt.Errorf("parse env vars: %w", err)
	}

	if z.args.ZoomToken == "" {
		return nil
	}

	if err := json.Unmarshal([]byte(z.args.ZoomToken), &z.token); err != nil {
		return fmt.Errorf("unmarshal token: %w", err)
	}

	return nil
}

func CreateMeeting(client *http.Client, topic, user string) (*http.Response, error) {
	params := zoom.CreateMeetingOptions{
		Type:  zoom.MeetingTypeInstant,
		Topic: topic,
	}

	var buf bytes.Buffer
	if err := json.NewEncoder(&buf).Encode(params); err != nil {
		return nil, fmt.Errorf("marshal params: %w", err)
	}

	resp, err := client.Post(fmt.Sprintf(base.ZoomAPIV2+zoom.CreateMeetingPath, user), "application/json", &buf)
	if err != nil {
		return nil, fmt.Errorf("post request: %w", err)
	}

	return resp, nil
}

func (z *ZoomMeetingCreate) Run() (exitCode int, output []byte, err error) {
	if z.args.ZoomApiKey != "" || z.args.ZoomApiSecret != "" {
		return z.DeprecatedRun()
	}

	if z.args.ZoomToken == "" {
		return step.ExitCodeFailure, nil, fmt.Errorf("zoom token not found")
	}

	client := oauth2.NewClient(context.Background(), oauth2.StaticTokenSource(z.token))

	resp, err := CreateMeeting(client, z.args.MeetingTopic, z.args.HostEmail)
	if err != nil {
		return step.ExitCodeFailure, nil, fmt.Errorf("create meeting: %w", err)
	}

	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return step.ExitCodeFailure, nil, fmt.Errorf("read response body: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		var apiError *zoom.APIError

		if err := json.Unmarshal(respBody, &apiError); err != nil {
			return step.ExitCodeFailure, nil, fmt.Errorf("unmarshal response: %w", err)
		}

		return step.ExitCodeFailure, respBody, apiError
	}

	return step.ExitCodeOK, respBody, nil
}

func (z *ZoomMeetingCreate) DeprecatedRun() (exitCode int, output []byte, err error) {
	zoom.APIKey = z.args.ZoomApiKey
	zoom.APISecret = z.args.ZoomApiSecret

	user, err := zoom.GetUser(zoom.GetUserOpts{EmailOrID: z.args.HostEmail})
	if err != nil {
		return step.ExitCodeFailure, nil, fmt.Errorf("failed to get zoom user: %w", err)
	}

	meeting, err := zoom.CreateMeeting(zoom.CreateMeetingOptions{
		Type:   zoom.MeetingTypeInstant,
		Topic:  z.args.MeetingTopic,
		HostID: user.ID,
	})
	if err != nil {
		return step.ExitCodeFailure, nil, err
	}

	marshaledMeeting, err := json.Marshal(meeting)
	if err != nil {
		return step.ExitCodeFailure, nil, err
	}

	return step.ExitCodeOK, marshaledMeeting, nil
}

func main() {
	step.Run(&ZoomMeetingCreate{})
}
