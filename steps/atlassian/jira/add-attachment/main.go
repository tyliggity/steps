package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"strings"

	"github.com/stackpulse/public-steps/atlassian/jira/base"
	"github.com/stackpulse/public-steps/common/env"
	"github.com/stackpulse/public-steps/common/step"
	"golang.org/x/oauth2"
)

type jiraAddAttachment struct {
	base.Args
	IssueID        string `env:"ISSUE_ID,required"`
	Attachment     string `env:"ATTACHMENT,required"`
	AttachmentName string `env:"ATTACHMENT_NAME,required"`
	token          *oauth2.Token
}

type output struct {
	Id   string `json:"id"`
	Self string `json:"self"`
	step.Outputs
}

func (j *jiraAddAttachment) Init() error {
	err := env.Parse(j)
	if err != nil {
		return err
	}

	if err := json.Unmarshal([]byte(j.JiraToken), &j.token); err != nil {
		return fmt.Errorf("unmarshal token: %w", err)
	}

	return nil
}

func (j *jiraAddAttachment) Run() (int, []byte, error) {
	jiraClient, err := base.NewClient(j.token, j.SiteName)
	if err != nil {
		return step.ExitCodeFailure, nil, err
	}

	attachments, resp, err := jiraClient.Issue.PostAttachment(j.IssueID, strings.NewReader(j.Attachment), j.AttachmentName)
	if err != nil {
		return step.ExitCodeFailure, base.ExtractResponse(resp), fmt.Errorf("post attachment to issue: %w", err)
	}

	if attachments == nil || len(*attachments) <= 0 {
		return step.ExitCodeFailure, base.ExtractResponse(resp), errors.New("no attachment")
	}

	attachment := (*attachments)[0]

	jsonOutput, err := json.Marshal(&output{
		Outputs: step.Outputs{Object: attachment},
		Id:      attachment.ID,
		Self:    attachment.Self,
	})
	if err != nil {
		return step.ExitCodeFailure, nil, fmt.Errorf("marshal output: %w", err)
	}

	return step.ExitCodeOK, jsonOutput, nil
}

func main() {
	step.Run(&jiraAddAttachment{})
}
