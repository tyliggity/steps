module github.com/stackpulse/public-steps/jenkins/last-build

go 1.14

require (
	github.com/davecgh/go-spew v1.1.1 // indirect
	github.com/kr/pretty v0.2.0 // indirect
	github.com/pkg/errors v0.9.1 // indirect
	github.com/stackpulse/public-steps/common v0.0.0
	github.com/stackpulse/public-steps/jenkins/base v0.0.0
	golang.org/x/oauth2 v0.0.0-20200902213428-5d25da1a8d43
	gopkg.in/check.v1 v1.0.0-20190902080502-41f04d3bba15 // indirect
)

replace github.com/stackpulse/public-steps/common v0.0.0 => ../../common
replace github.com/stackpulse/public-steps/jenkins/base v0.0.0 => ../base
