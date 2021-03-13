module github.com/stackpulse/public-steps/opsgenie/alert/acknowledge

go 1.14

replace (
	github.com/stackpulse/public-steps/common v0.0.0 => ../../../common
	github.com/stackpulse/public-steps/opsgenie/base v0.0.0 => ../../base
)

require (
	github.com/Jeffail/gabs/v2 v2.6.0
	github.com/caarlos0/env/v6 v6.5.0
	github.com/opsgenie/opsgenie-go-sdk-v2 v1.2.1
	github.com/pkg/errors v0.9.1 // indirect
	github.com/sirupsen/logrus v1.4.2
	github.com/stackpulse/public-steps/common v0.0.0
	github.com/stackpulse/public-steps/opsgenie/base v0.0.0
)
