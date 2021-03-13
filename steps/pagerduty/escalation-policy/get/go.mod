module github.com/stackpulse/public-steps/pagerduty/escalation-policy/get

go 1.14

replace (
	github.com/stackpulse/public-steps/common v0.0.0 => ../../../common
	github.com/stackpulse/public-steps/pagerduty/base v0.0.0 => ../../base
)

require (
	github.com/PagerDuty/go-pagerduty v1.3.0
	github.com/stackpulse/public-steps/common v0.0.0
	github.com/stackpulse/public-steps/pagerduty/base v0.0.0
)
