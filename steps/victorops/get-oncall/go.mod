module github.com/stackpulse/public-steps/victorops/get-oncall

go 1.14

replace (
	github.com/stackpulse/public-steps/common v0.0.0 => ../../common
	github.com/stackpulse/public-steps/victorops/base v0.0.0 => ../base
)

require (
	github.com/caarlos0/env/v6 v6.5.0
	github.com/stackpulse/public-steps/common v0.0.0
	github.com/stackpulse/public-steps/victorops/base v0.0.0
	github.com/victorops/go-victorops v1.0.1 // indirect
)
