module github.com/stackpulse/steps/psql/query/terminate

go 1.14

require (
	github.com/Jeffail/gabs v1.4.0
	github.com/lib/pq v1.9.0
	github.com/stackpulse/steps-sdk-go v0.0.0-20210314133745-61086c27983f
	github.com/stackpulse/steps/psql/base v0.0.0
	github.com/stackpulse/steps/psql/query/base v0.0.0
)

replace (
	github.com/stackpulse/steps/psql/base v0.0.0 => ../../base
	github.com/stackpulse/steps/psql/query/base v0.0.0 => ../base
	github.com/stackpulse/steps/utils/base/env v0.0.0 => ../../../utils/base/env
)
