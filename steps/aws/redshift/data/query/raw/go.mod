module github.com/stackpulse/steps/steps/redshift/data/query/raw

go 1.16

require (
	github.com/stackpulse/steps-sdk-go v0.0.0-20210329111118-cc87e9772586
	github.com/stackpulse/steps/redshift/base v0.0.0
)

replace github.com/stackpulse/steps/redshift/base v0.0.0 => ./../../../base
