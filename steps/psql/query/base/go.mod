module github.com/stackpulse/public-steps/psql/query/base

go 1.14

require (
    github.com/stackpulse/public-steps/psql/base v0.0.0
    github.com/stackpulse/public-steps/common v0.0.0
)

replace (
	github.com/stackpulse/public-steps/common v0.0.0 => ../../../common
	github.com/stackpulse/public-steps/psql/base v0.0.0 => ../../base
)
