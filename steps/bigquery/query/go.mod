module github.com/stackpulse/steps/bigquery/query

go 1.14

require (
	cloud.google.com/go/bigquery v1.12.0
	github.com/caarlos0/env/v6 v6.5.0
	github.com/stackpulse/steps/common v0.0.0
	google.golang.org/api v0.32.0
)

replace github.com/stackpulse/steps/common v0.0.0 => ../../common
