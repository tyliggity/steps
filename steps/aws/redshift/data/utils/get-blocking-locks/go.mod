module github.com/stackpulse/steps/aws/redshift/data/amazon-redshift-utils/get-blocking-locks

go 1.16

require (
	github.com/stackpulse/steps-sdk-go v0.0.0-20210329111118-cc87e9772586
	github.com/stackpulse/steps/aws/redshift/base v0.0.0
	github.com/stretchr/testify v1.7.0
)

replace github.com/stackpulse/steps/aws/redshift/base v0.0.0 => ./../../../base
