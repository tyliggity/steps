module github.com/stackpulse/steps/steps/kubectl/cmd

go 1.14

require (
	github.com/stackpulse/steps-sdk-go v0.0.0-20210314133745-61086c27983f
	github.com/stackpulse/steps/kubectl/base v0.0.0
)

replace github.com/stackpulse/steps/kubectl/base v0.0.0 => ../base
