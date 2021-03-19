module github.com/stackpulse/public-steps/rabbitmq/overview

go 1.14

require (
	github.com/stackpulse/public-steps/rabbitmq/base v0.0.0
	github.com/stackpulse/steps-sdk-go v0.0.0-20210314133745-61086c27983f
)

replace (
	github.com/stackpulse/public-steps/rabbitmq/base v0.0.0 => ../base
	github.com/stackpulse/steps-sdk-go v0.0.0-20210314133745-61086c27983f => ../../common
)
