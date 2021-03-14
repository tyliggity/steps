module github.com/stackpulse/public-steps/pagerduty/incident/list

go 1.14

replace (
	github.com/stackpulse/public-steps/common v0.0.0 => ../../../common
	github.com/stackpulse/public-steps/pagerduty/base v0.0.0 => ../../base
	github.com/stackpulse/public-steps/pagerduty/incident/base v0.0.0 => ../base
)

require (
	github.com/PagerDuty/go-pagerduty v1.3.0
	github.com/caarlos0/env/v6 v6.5.0
	github.com/stackpulse/public-steps/common v0.0.0
	github.com/stackpulse/public-steps/pagerduty/base v0.0.0
	github.com/stackpulse/public-steps/pagerduty/incident/base v0.0.0
)
