module github.com/stackpulse/public-steps/steps/ansible/awx/get-job

go 1.15

require (
	github.com/stackpulse/public-steps/steps/ansible/awx/base v0.0.0
	github.com/stackpulse/steps-sdk-go v0.0.0-20210314133745-61086c27983f
)

replace github.com/stackpulse/public-steps/steps/ansible/awx/base v0.0.0 => ../base
