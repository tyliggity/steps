module github.com/stackpulse/public-steps/aws/redshift/data/query/cancel

go 1.16

require (
	cloud.google.com/go/storage v1.14.0 // indirect
	github.com/stackpulse/public-steps/aws/redshift/base v0.0.0
	github.com/stackpulse/public-steps/common v0.0.0
)

replace (
	github.com/stackpulse/public-steps/aws/redshift/base v0.0.0 => ./../../../base
	github.com/stackpulse/public-steps/common v0.0.0 => ./../../../../../common
)
