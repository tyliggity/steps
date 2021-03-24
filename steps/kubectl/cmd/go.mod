module github.com/stackpulse/public-steps/steps/kubectl/cmd

go 1.14

require (
	github.com/stackpulse/public-steps/kubectl/base v0.0.0
	github.com/stackpulse/steps-sdk-go v0.0.0-20210314133745-61086c27983f
)

replace github.com/stackpulse/public-steps/kubectl/base v0.0.0 => ../base
