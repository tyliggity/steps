module github.com/stackpulse/public-steps/aws/redshift/data/amazon-redshift-utils/get-blocking-locks

go 1.16

require (
	github.com/stackpulse/public-steps/aws/redshift/base v0.0.0
	github.com/stackpulse/steps-sdk-go v0.0.0-20210314133745-61086c27983f
	github.com/stretchr/testify v1.7.0
)

replace (
	github.com/stackpulse/public-steps/aws/redshift/base v0.0.0 => ./../../../base
	github.com/stackpulse/steps-sdk-go v0.0.0-20210314133745-61086c27983f => ./../../../../../common
)
