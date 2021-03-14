module github.com/stackpulse/public-steps/psql/query/raw

go 1.14

require (
    github.com/stackpulse/public-steps/psql/base v0.0.0
    github.com/stackpulse/public-steps/psql/query/base v0.0.0
    github.com/stackpulse/public-steps/common v0.0.0
)

replace (
	github.com/stackpulse/public-steps/common v0.0.0 => ../../../common
    github.com/stackpulse/public-steps/psql/query/base v0.0.0 => ../base
	github.com/stackpulse/public-steps/psql/base v0.0.0 => ../../base
	github.com/stackpulse/public-steps/utils/base/env v0.0.0 => ../../../utils/base/env
)
