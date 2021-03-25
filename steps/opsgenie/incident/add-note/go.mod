module github.com/stackpulse/steps/opsgenie/incident/add-note

go 1.14

replace github.com/stackpulse/steps/opsgenie/base v0.0.0 => ../../base

require (
	github.com/Jeffail/gabs/v2 v2.6.0
	github.com/caarlos0/env/v6 v6.5.0
	github.com/opsgenie/opsgenie-go-sdk-v2 v1.2.6
	github.com/pkg/errors v0.9.1 // indirect
	github.com/stackpulse/steps-sdk-go v0.0.0-20210314133745-61086c27983f
	github.com/stackpulse/steps/opsgenie/base v0.0.0
)
