module github.com/stackpulse/public-steps/istio/proxy-status

go 1.14

require (
	github.com/caarlos0/env/v6 v6.3.0
	github.com/stackpulse/public-steps/common v0.0.0
	github.com/stackpulse/public-steps/istio/base v0.0.0
)

replace github.com/stackpulse/public-steps/common v0.0.0 => ../../common
replace github.com/stackpulse/public-steps/istio/base v0.0.0 => ../base
