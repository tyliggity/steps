module github.com/stackpulse/steps/redis/hget

go 1.14

require (
	github.com/Jeffail/gabs/v2 v2.6.0
	github.com/caarlos0/env/v6 v6.5.0
	github.com/go-redis/redis/v8 v8.4.4
	github.com/stackpulse/steps/common v0.0.0
	github.com/stackpulse/steps/redis/base v0.0.0
)

replace github.com/stackpulse/steps/common v0.0.0 => ../../common

replace github.com/stackpulse/steps/redis/base v0.0.0 => ../base
