module github.com/stackpulse/public-steps/jenkins/last-build

go 1.14

require (
	github.com/davecgh/go-spew v1.1.1 // indirect
	github.com/kr/pretty v0.2.0 // indirect
	github.com/stackpulse/public-steps/jenkins/base v0.0.0
	github.com/stackpulse/steps-sdk-go v0.0.0-20210314133745-61086c27983f
	golang.org/x/net v0.0.0-20200822124328-c89045814202 // indirect
	gopkg.in/check.v1 v1.0.0-20190902080502-41f04d3bba15 // indirect
)

replace github.com/stackpulse/public-steps/jenkins/base v0.0.0 => ../base
