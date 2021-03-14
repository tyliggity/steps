module github.com/stackpulse/public-steps/logzio

go 1.14

require (
	github.com/caarlos0/env/v6 v6.5.0
	github.com/stackpulse/public-steps/common v0.0.0
	github.com/stackpulse/public-steps/logzio/base v0.0.0
)

replace github.com/stackpulse/public-steps/common v0.0.0 => ../../common

replace github.com/stackpulse/public-steps/logzio/base v0.0.0 => ../base
