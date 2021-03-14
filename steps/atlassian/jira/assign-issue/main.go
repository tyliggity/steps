package main

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/stackpulse/steps/atlassian/jira/base"
	"github.com/stackpulse/steps/common/env"
	"github.com/stackpulse/steps/common/step"
	"golang.org/x/oauth2"
)

type jiraAssignIssue struct {
	base.Args
	IssueID       string `env:"ISSUE_ID,required"`
	AssigneeEmail string `env:"ASSIGNEE_EMAIL,required"`
	token         *oauth2.Token
}

type output struct {
	AssigneeId string `json:"assignee_id"`
}

func (j *jiraAssignIssue) Init() error {
	err := env.Parse(j)
	if err != nil {
		return err
	}

	if err := json.Unmarshal([]byte(j.JiraToken), &j.token); err != nil {
		return fmt.Errorf("unmarshal token: %w", err)
	}

	return nil
}

func (j *jiraAssignIssue) Run() (int, []byte, error) {
	jiraClient, err := base.NewClient(j.token, j.SiteName)
	if err != nil {
		return step.ExitCodeFailure, nil, err
	}

	results, resp, err := jiraClient.User.Find(j.AssigneeEmail)
	if err != nil {
		return step.ExitCodeFailure, base.ExtractResponse(resp), fmt.Errorf("find user by email: %w", err)
	}

	if results == nil || len(results) < 1 {
		return step.ExitCodeFailure, base.ExtractResponse(resp), fmt.Errorf("assignee not found")
	}

	assignee := &results[0]

	resp, err = jiraClient.Issue.UpdateAssignee(j.IssueID, assignee)
	if err != nil {
		return step.ExitCodeFailure, base.ExtractResponse(resp), fmt.Errorf("update issue assignee: %w", err)
	}

	if resp.StatusCode != http.StatusNoContent {
		return step.ExitCodeFailure, base.ExtractResponse(resp), fmt.Errorf("invalid status code: %d expected: %d", resp.StatusCode, http.StatusNoContent)
	}

	jsonOutput, err := json.Marshal(output{
		AssigneeId: assignee.AccountID,
	})
	if err != nil {
		return step.ExitCodeFailure, base.ExtractResponse(resp), fmt.Errorf("marshal output: %w", err)
	}

	return step.ExitCodeOK, jsonOutput, nil
}

func main() {
	step.Run(&jiraAssignIssue{})
}
