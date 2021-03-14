module github.com/stackpulse/steps/prometheus/query/raw

go 1.14

require (
    github.com/stackpulse/steps/prometheus/base v0.0.0
    github.com/stackpulse/steps/prometheus/query/base v0.0.0
)

replace (
    github.com/stackpulse/steps/prometheus/query/base v0.0.0 => ../base
	github.com/stackpulse/steps/prometheus/base v0.0.0 => ../../base
)
