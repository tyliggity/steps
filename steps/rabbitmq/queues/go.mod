module github.com/stackpulse/steps/rabbitmq/queues

go 1.14

require (
	github.com/Jeffail/gabs/v2 v2.6.0
	github.com/caarlos0/env/v6 v6.5.0
	github.com/stackpulse/steps-sdk-go v0.0.0-20210314133745-61086c27983f
	github.com/stackpulse/steps/rabbitmq/base v0.0.0
)

replace github.com/stackpulse/steps/rabbitmq/base v0.0.0 => ../base
