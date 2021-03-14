package main

import (
	"encoding/json"
	"fmt"

	"github.com/andygrunwald/go-jira"
	"github.com/stackpulse/steps/atlassian/jira/base"
	envconf "github.com/stackpulse/steps/common/env"
	"github.com/stackpulse/steps/common/step"
	"golang.org/x/oauth2"
)

type Args struct {
	base.Args
	IssueID string `env:"ISSUE_ID,required"`
	Comment string `env:"COMMENT,required"`
}

type output struct {
	Id     string `json:"id"`
	Name   string `json:"name"`
	Author string `json:"author"`
	step.Outputs
}

type jiraAddComment struct {
	args  Args
	token *oauth2.Token
}

func (j *jiraAddComment) Init() error {
	err := envconf.Parse(&j.args)
	if err != nil {
		return err
	}

	if err := json.Unmarshal([]byte(j.args.JiraToken), &j.token); err != nil {
		return fmt.Errorf("unmarshal token: %w", err)
	}

	return nil
}

func (j *jiraAddComment) Run() (int, []byte, error) {
	oauthClient, url, err := base.GetOauthConnectionDetails(j.token, j.args.SiteName)
	if err != nil {
		return step.ExitCodeFailure, nil, err
	}

	jiraClient, err := jira.NewClient(oauthClient, url)
	if err != nil {
		return step.ExitCodeFailure, nil, err
	}

	var resp *jira.Response
	newComment, resp, err := jiraClient.Issue.AddComment(j.args.IssueID, &jira.Comment{Body: j.args.Comment})
	if err != nil {
		return step.ExitCodeFailure, base.ExtractResponse(resp), err
	}

	jsonOutput, err := json.Marshal(&output{
		Outputs: step.Outputs{Object: newComment},
		Id:      newComment.ID,
		Name:    newComment.Name,
		Author:  newComment.Author.Name,
	})
	if err != nil {
		return step.ExitCodeFailure, nil, err
	}

	return step.ExitCodeOK, jsonOutput, nil
}

func main() {
	step.Run(&jiraAddComment{})
}
