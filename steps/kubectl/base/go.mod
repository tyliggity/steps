module github.com/stackpulse/public-steps/kubectl/base

go 1.14

require (
	github.com/Jeffail/gabs/v2 v2.6.0
	github.com/caarlos0/env/v6 v6.5.0
	github.com/google/shlex v0.0.0-20191202100458-e7afc7fbc510 // indirect
	github.com/hashicorp/go-multierror v1.1.0
	github.com/stackpulse/public-steps/common v0.0.0
	maze.io/x/duration.v1 v0.0.0-20161004121933-0bd39bea6019
)

replace github.com/stackpulse/public-steps/common v0.0.0 => ../../common
