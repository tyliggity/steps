module github.com/stackpulse/steps/steps/kubectl/jattach-heapdump

go 1.16

require (
	github.com/google/shlex v0.0.0-20191202100458-e7afc7fbc510
	github.com/stackpulse/steps-sdk-go v0.0.0-20210329111118-cc87e9772586
	github.com/stackpulse/steps/kubectl/base v0.0.0
	gocloud.dev v0.22.0
)

replace github.com/stackpulse/steps/kubectl/base v0.0.0 => ../base
