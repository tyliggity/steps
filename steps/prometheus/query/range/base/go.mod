module github.com/stackpulse/steps/prometheus/query/range/base

go 1.14

require (
	github.com/caarlos0/env/v6 v6.3.0
	github.com/prometheus/common v0.10.0
	github.com/stackpulse/steps/prometheus/query/base v0.0.0
)

replace (
	github.com/stackpulse/steps/prometheus/base v0.0.0 => ../../../base
	github.com/stackpulse/steps/prometheus/query/base v0.0.0 => ../../base
)
