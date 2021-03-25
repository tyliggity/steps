module github.com/stackpulse/steps/steps/ansible/awx/get-job

go 1.15

require (
	github.com/stackpulse/steps-sdk-go v0.0.0-20210314133745-61086c27983f
	github.com/stackpulse/steps/steps/ansible/awx/base v0.0.0
)

replace github.com/stackpulse/steps/steps/ansible/awx/base v0.0.0 => ../base
