module github.com/stackpulse/steps/pagerduty/incident/base

go 1.14

replace github.com/stackpulse/steps/pagerduty/base v0.0.0 => ../../base

require (
	github.com/PagerDuty/go-pagerduty v1.3.0
	github.com/stackpulse/steps/pagerduty/base v0.0.0
)
