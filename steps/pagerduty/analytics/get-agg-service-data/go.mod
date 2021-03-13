module github.com/stackpulse/public-steps/pagerduty/analytics/get-agg-service-data

go 1.15

replace (
	github.com/stackpulse/public-steps/common v0.0.0 => ../../../common
	github.com/stackpulse/public-steps/pagerduty/base v0.0.0 => ../../base
)

require (
	github.com/stackpulse/public-steps/common v0.0.0
	github.com/stackpulse/public-steps/pagerduty/base v0.0.0
)
