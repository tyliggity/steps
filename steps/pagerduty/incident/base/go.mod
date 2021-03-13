module github.com/stackpulse/steps/pagerduty/incident/base

go 1.14

replace (
	github.com/stackpulse/steps/common v0.0.0 => ../../../common
	github.com/stackpulse/steps/pagerduty/base v0.0.0 => ../../base
)

require (
	github.com/PagerDuty/go-pagerduty v1.3.0
	github.com/caarlos0/env/v6 v6.5.0
	github.com/pkg/errors v0.9.1
	github.com/stackpulse/steps/common v0.0.0
	github.com/stackpulse/steps/pagerduty/base v0.0.0
)
