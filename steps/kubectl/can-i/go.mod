module github.com/stackpulse/public-steps/steps/kubectl/can_i

go 1.14

require (
	github.com/stackpulse/public-steps/common v0.0.0
	github.com/stackpulse/public-steps/kubectl/base v0.0.0 // indirect
)

replace github.com/stackpulse/public-steps/common v0.0.0 => ../../common

replace github.com/stackpulse/public-steps/kubectl/base v0.0.0 => ../base
