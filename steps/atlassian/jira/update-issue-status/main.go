package main

import (
	"encoding/json"
	"fmt"
	"github.com/stackpulse/steps/common/env"

	"github.com/andygrunwald/go-jira"
	"github.com/stackpulse/steps/atlassian/jira/base"
	"github.com/stackpulse/steps/common/step"
	"golang.org/x/oauth2"
)

type Args struct {
	base.Args
	IssueID        string `env:"ISSUE_ID,required"`
	NewIssueStatus string `env:"NEW_ISSUE_STATUS,required"`
}

type JiraUpdateIssueStatus struct {
	args  Args
	token *oauth2.Token
}

type Outputs struct {
	Success bool `json:"success"`
	step.Outputs
}

func (j *JiraUpdateIssueStatus) Init() error {
	err := env.Parse(&j.args)
	if err != nil {
		return err
	}

	if err := json.Unmarshal([]byte(j.args.JiraToken), &j.token); err != nil {
		return fmt.Errorf("unmarshal token: %w", err)
	}

	return nil
}

func (j *JiraUpdateIssueStatus) Run() (int, []byte, error) {
	oauthClient, url, err := base.GetOauthConnectionDetails(j.token, j.args.SiteName)
	if err != nil {
		return step.ExitCodeFailure, nil, err
	}

	jiraClient, err := jira.NewClient(oauthClient, url)
	if err != nil {
		return step.ExitCodeFailure, nil, err
	}

	possibleTransitions, _, err := jiraClient.Issue.GetTransitions(j.args.IssueID)
	if err != nil {
		return step.ExitCodeFailure, nil, fmt.Errorf("could't get issue transitions: %w", err)
	}

	var transitionID string
	for _, v := range possibleTransitions {
		if v.Name == j.args.NewIssueStatus {
			transitionID = v.ID
			break
		}
	}

	if transitionID == "" {
		return step.ExitCodeFailure, nil, fmt.Errorf("couldn't find transition id for issue status: %s "+
			"check that your account has that issue status", j.args.NewIssueStatus)
	}

	_, err = jiraClient.Issue.DoTransition(j.args.IssueID, transitionID)
	if err != nil {
		return step.ExitCodeFailure, nil, err
	}

	marshaledOutput, err := json.Marshal(&Outputs{Success: true})
	if err != nil {
		return step.ExitCodeFailure, nil, err
	}

	return step.ExitCodeOK, marshaledOutput, nil
}

func main() {
	step.Run(&JiraUpdateIssueStatus{})
}
