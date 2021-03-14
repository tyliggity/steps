module github.com/stackpulse/public-steps/steps/redshift/data/query/raw

go 1.16

require (
	github.com/stackpulse/public-steps/common v0.0.0
	github.com/stackpulse/public-steps/redshift/base v0.0.0
)

replace (
	github.com/stackpulse/public-steps/common v0.0.0 => ./../../../../../common
	github.com/stackpulse/public-steps/redshift/base v0.0.0 => ./../../../base
)
