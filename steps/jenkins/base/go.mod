module github.com/stackpulse/public-steps/steps/jenkins/base

go 1.15

require (
	github.com/bndr/gojenkins v1.0.2-0.20210112054307-ab81397930ca
	github.com/stackpulse/steps-sdk-go v0.0.0-20210314133745-61086c27983f
)

replace github.com/stackpulse/steps-sdk-go v0.0.0-20210314133745-61086c27983f => ../../common
