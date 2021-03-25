module github.com/stackpulse/steps/atlassian/jira/delete-issue

go 1.14

require (
	github.com/davecgh/go-spew v1.1.1 // indirect
	github.com/kr/pretty v0.2.0 // indirect
	github.com/pkg/errors v0.9.1 // indirect
	github.com/stackpulse/steps-sdk-go v0.0.0-20210314133745-61086c27983f
	github.com/stackpulse/steps/atlassian/jira/base v0.0.0
	golang.org/x/oauth2 v0.0.0-20200902213428-5d25da1a8d43
	gopkg.in/check.v1 v1.0.0-20190902080502-41f04d3bba15 // indirect
)

replace github.com/stackpulse/steps/atlassian/jira/base v0.0.0 => ../base
