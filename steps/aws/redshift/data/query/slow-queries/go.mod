module github.com/stackpulse/public-steps/steps/redshift/data/query/slow-queries

go 1.16

require (
	github.com/stackpulse/public-steps/common v0.0.0
	github.com/stackpulse/public-steps/redshift/base v0.0.0
	github.com/stretchr/testify v1.7.0
)

replace (
	github.com/stackpulse/public-steps/common v0.0.0 => ./../../../../../common
	github.com/stackpulse/public-steps/redshift/base v0.0.0 => ./../../../base
)
