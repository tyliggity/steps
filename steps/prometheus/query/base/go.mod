module github.com/stackpulse/public-steps/prometheus/query/base

go 1.14

require (
    github.com/stackpulse/public-steps/prometheus/base v0.0.0
)

replace (
	github.com/stackpulse/public-steps/prometheus/base v0.0.0 => ../../base
)
