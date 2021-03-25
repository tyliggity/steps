module github.com/stackpulse/steps/steps/elastic/cluster-health

go 1.14

replace github.com/stackpulse/steps/elastic/base v0.0.0 => ../base

require (
	github.com/caarlos0/env/v6 v6.3.0
	github.com/stackpulse/steps/elastic/base v0.0.0
)
