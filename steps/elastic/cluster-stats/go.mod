module github.com/stackpulse/public-steps/steps/elastic/cluster-stats

go 1.14

replace (
	github.com/stackpulse/public-steps/elastic/base v0.0.0 => ../base
	github.com/stackpulse/steps-sdk-go v0.0.0-20210314133745-61086c27983f => ../../common
)

require (
	github.com/caarlos0/env/v6 v6.3.0
	github.com/stackpulse/public-steps/elastic/base v0.0.0
)
