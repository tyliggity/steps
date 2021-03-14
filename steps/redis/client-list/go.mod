module github.com/stackpulse/public-steps/redis/client-list

go 1.14

require (
	github.com/caarlos0/env/v6 v6.5.0
	github.com/fzzy/radix v0.5.6
	github.com/go-redis/redis/v8 v8.4.4
	github.com/stackpulse/public-steps/common v0.0.0
	github.com/stackpulse/public-steps/redis/base v0.0.0
)

replace github.com/stackpulse/public-steps/common v0.0.0 => ../../common

replace github.com/stackpulse/public-steps/redis/base v0.0.0 => ../base
