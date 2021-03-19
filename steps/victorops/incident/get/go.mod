module github.com/stackpulse/public-steps/victorops/incident/get

go 1.14

require (
	github.com/davecgh/go-spew v1.1.1 // indirect
	github.com/kr/pretty v0.2.0 // indirect
	github.com/stackpulse/public-steps/victorops/base v0.0.0
	github.com/stackpulse/steps-sdk-go v0.0.0-20210314133745-61086c27983f
	github.com/victorops/go-victorops v1.0.4
	gopkg.in/check.v1 v1.0.0-20190902080502-41f04d3bba15 // indirect
)

replace github.com/stackpulse/public-steps/victorops/base v0.0.0 => ../../base
