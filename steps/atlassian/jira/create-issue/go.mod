module github.com/stackpulse/steps/steps/atlassian/jira/create-issue

go 1.14

require (
	github.com/andygrunwald/go-jira v1.13.0
	github.com/stackpulse/steps-sdk-go v0.0.0-20210314133745-61086c27983f
	github.com/stackpulse/steps/atlassian/jira/base v0.0.0
	golang.org/x/oauth2 v0.0.0-20200902213428-5d25da1a8d43
)

replace github.com/stackpulse/steps/atlassian/jira/base v0.0.0 => ../base
