module github.com/stackpulse/public-steps/zoom/meeting/create

go 1.14

replace (
	github.com/stackpulse/public-steps/common v0.0.0 => ../../../common
	github.com/stackpulse/public-steps/zoom/base v0.0.0 => ../../base
)

require (
	github.com/caarlos0/env/v6 v6.5.0
	github.com/davecgh/go-spew v1.1.1 // indirect
	github.com/himalayan-institute/zoom-lib-golang v1.0.0
	github.com/stackpulse/public-steps/common v0.0.0
	github.com/stackpulse/public-steps/zoom/base v0.0.0
)
