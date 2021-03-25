module github.com/stackpulse/steps/rabbitmq/consumers

go 1.14

require (
	github.com/stackpulse/steps-sdk-go v0.0.0-20210314133745-61086c27983f
	github.com/stackpulse/steps/rabbitmq/base v0.0.0
)

replace github.com/stackpulse/steps/rabbitmq/base v0.0.0 => ../base
