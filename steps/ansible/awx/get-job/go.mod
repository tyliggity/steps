module github.com/stackpulse/public-steps/public-steps/ansible/awx/get-job

go 1.15

require (
	github.com/stackpulse/public-steps/common v0.0.0
	github.com/stackpulse/public-steps/public-steps/ansible/awx/base v0.0.0
)

replace (
 	github.com/stackpulse/public-steps/public-steps/ansible/awx/base v0.0.0 => ../base
 	github.com/stackpulse/public-steps/common v0.0.0 => ../../../common
)