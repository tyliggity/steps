module github.com/stackpulse/steps/coralogix/query

go 1.14

require (
	github.com/caarlos0/env/v6 v6.5.0
	github.com/stackpulse/steps/common v0.0.0
	github.com/stackpulse/steps/coralogix/base v0.0.0
)

replace github.com/stackpulse/steps/common v0.0.0 => ../../common

replace github.com/stackpulse/steps/coralogix/base v0.0.0 => ../base
