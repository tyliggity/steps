module github.com/stackpulse/public-steps/opsgenie/oncall/get

go 1.14

replace github.com/stackpulse/public-steps/opsgenie/base v0.0.0 => ../../base

require (
	github.com/caarlos0/env/v6 v6.5.0
	github.com/opsgenie/opsgenie-go-sdk-v2 v1.2.6
	github.com/pkg/errors v0.9.1 // indirect
	github.com/stackpulse/public-steps/opsgenie/base v0.0.0
	github.com/stackpulse/steps-sdk-go v0.0.0-20210314133745-61086c27983f
)
