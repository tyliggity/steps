module github.com/stackpulse/public-steps/steps/jenkins/run

go 1.16

require (
	github.com/stackpulse/public-steps/steps/jenkins/base v0.0.0
	github.com/stackpulse/steps-sdk-go v0.0.0-20210314133745-61086c27983f
)

replace (
	github.com/stackpulse/public-steps/steps/jenkins/base v0.0.0 => ../base
	github.com/stackpulse/steps-sdk-go v0.0.0-20210314133745-61086c27983f => ../../common
)
