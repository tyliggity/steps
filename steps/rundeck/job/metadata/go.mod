module github.com/stackpulse/public-steps/public-steps/rundeck/job/metadata

go 1.15

require (
	github.com/stackpulse/public-steps/common v0.0.0
	github.com/stackpulse/public-steps/public-steps/rundeck/base v0.0.0
)

replace (
 	github.com/stackpulse/public-steps/public-steps/rundeck/base v0.0.0 => ../../base
 	github.com/stackpulse/public-steps/common v0.0.0 => ../../../common
)