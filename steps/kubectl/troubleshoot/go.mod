module github.com/stackpulse/steps/steps/kubectl/troubleshoot

go 1.14

require (
	github.com/stackpulse/steps/kubectl/base v0.0.0
	github.com/olekukonko/tablewriter v0.0.5
	github.com/stackpulse/steps-sdk-go v0.0.0-20210314133745-61086c27983f
)

replace github.com/stackpulse/steps/kubectl/base v0.0.0 => ../base
