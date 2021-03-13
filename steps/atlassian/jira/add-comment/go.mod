module github.com/stackpulse/public-steps/public-steps/atlassian/jira/add-comment

go 1.14

require (
	github.com/andygrunwald/go-jira v1.13.0
	github.com/stackpulse/public-steps/atlassian/jira/base v0.0.0
	github.com/stackpulse/public-steps/common v0.0.0
	golang.org/x/oauth2 v0.0.0-20200902213428-5d25da1a8d43
)

replace github.com/stackpulse/public-steps/common v0.0.0 => ../../../common

replace github.com/stackpulse/public-steps/atlassian/jira/base v0.0.0 => ../base
