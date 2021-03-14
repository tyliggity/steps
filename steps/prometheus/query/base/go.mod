module github.com/stackpulse/steps/prometheus/query/base

go 1.14

require (
    github.com/stackpulse/steps/prometheus/base v0.0.0
)

replace (
	github.com/stackpulse/steps/prometheus/base v0.0.0 => ../../base
)
