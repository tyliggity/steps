module github.com/stackpulse/steps/prometheus/query/base

go 1.14

require (
	github.com/prometheus/client_golang v1.7.1
	github.com/prometheus/common v0.10.0
	github.com/stackpulse/steps/prometheus/base v0.0.0
)

replace github.com/stackpulse/steps/prometheus/base v0.0.0 => ../../base
