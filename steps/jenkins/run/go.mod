module github.com/stackpulse/public-steps/public-steps/jenkins/run

go 1.16

require (
	github.com/bndr/gojenkins v1.0.2-0.20210112054307-ab81397930ca // indirect
	github.com/caarlos0/env/v6 v6.5.0 // indirect
	github.com/stackpulse/public-steps/common v0.0.0
	github.com/stackpulse/public-steps/public-steps/jenkins/base v0.0.0
)

replace (
	github.com/stackpulse/public-steps/common v0.0.0 => ../../common
	github.com/stackpulse/public-steps/public-steps/jenkins/base v0.0.0 => ../base
)
