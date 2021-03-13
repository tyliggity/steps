module github.com/stackpulse/public-steps/pagerduty/incident/acknowledge

go 1.14

replace (
	github.com/stackpulse/public-steps/common v0.0.0 => ../../../common
	github.com/stackpulse/public-steps/pagerduty/base v0.0.0 => ../../base
	github.com/stackpulse/public-steps/pagerduty/incident/base v0.0.0 => ../base
)

require (
	github.com/Jeffail/gabs/v2 v2.6.0
	github.com/stackpulse/public-steps/common v0.0.0
	github.com/stackpulse/public-steps/pagerduty/incident/base v0.0.0
)
