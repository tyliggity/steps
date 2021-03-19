module github.com/stackpulse/public-steps/aws/command

go 1.14

require (
	github.com/aws/aws-sdk-go-v2/config v1.0.0
	github.com/aws/aws-sdk-go-v2/service/ec2 v1.0.0
	github.com/caarlos0/env/v6 v6.5.0
	github.com/stackpulse/steps-sdk-go v0.0.0-20210314133745-61086c27983f
)

replace github.com/stackpulse/steps-sdk-go v0.0.0-20210314133745-61086c27983f => ../../common
