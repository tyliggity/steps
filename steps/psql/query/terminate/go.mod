module github.com/stackpulse/public-steps/psql/query/terminate

go 1.14

require (
	github.com/Jeffail/gabs v1.4.0
	github.com/lib/pq v1.9.0
	github.com/stackpulse/public-steps/psql/base v0.0.0
	github.com/stackpulse/public-steps/psql/query/base v0.0.0
	github.com/stackpulse/steps-sdk-go v0.0.0-20210314133745-61086c27983f
)

replace (
	github.com/stackpulse/public-steps/psql/base v0.0.0 => ../../base
	github.com/stackpulse/public-steps/psql/query/base v0.0.0 => ../base
	github.com/stackpulse/public-steps/utils/base/env v0.0.0 => ../../../utils/base/env
)
