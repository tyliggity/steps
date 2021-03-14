module github.com/stackpulse/public-steps/steps/jenkins/base

go 1.15

require (
	github.com/bndr/gojenkins v1.0.2-0.20210112054307-ab81397930ca // indirect
	github.com/stackpulse/public-steps/common v0.0.0
)

replace github.com/stackpulse/public-steps/common v0.0.0 => ../../common
