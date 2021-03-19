module github.com/stackpulse/public-steps/zoom/meeting/create

go 1.14

replace github.com/stackpulse/public-steps/zoom/base v0.0.0 => ../../base

require (
	github.com/caarlos0/env/v6 v6.5.0
	github.com/davecgh/go-spew v1.1.1 // indirect
	github.com/himalayan-institute/zoom-lib-golang v1.0.0
	github.com/stackpulse/public-steps/zoom/base v0.0.0
	github.com/stackpulse/steps-sdk-go v0.0.0-20210314133745-61086c27983f
	golang.org/x/oauth2 v0.0.0-20200107190931-bf48bf16ab8d
)
