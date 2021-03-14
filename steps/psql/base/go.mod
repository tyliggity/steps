module github.com/stackpulse/public-steps/psql/base

go 1.14

require (
	github.com/Jeffail/gabs/v2 v2.6.0
	github.com/caarlos0/env/v6 v6.5.0
	github.com/stackpulse/public-steps/common v0.0.0
)

replace github.com/stackpulse/public-steps/common v0.0.0 => ../../common
