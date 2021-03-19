module github.com/stackpulse/public-steps/pagerduty/incident/create

go 1.14

replace (
	github.com/stackpulse/public-steps/pagerduty/base v0.0.0 => ../../base
	github.com/stackpulse/public-steps/pagerduty/incident/base v0.0.0 => ../base
)

require (
	github.com/PagerDuty/go-pagerduty v1.3.0
	github.com/stackpulse/public-steps/pagerduty/incident/base v0.0.0
	github.com/stackpulse/steps-sdk-go v0.0.0-20210314133745-61086c27983f
)
