module github.com/stackpulse/steps/flux/lock

go 1.15

require (
	github.com/fluxcd/flux v1.17.0
	github.com/opencontainers/go-digest v1.0.0 // indirect
	github.com/spf13/pflag v1.0.5 // indirect
	github.com/stackpulse/steps-sdk-go v0.0.0-20210314133745-61086c27983f
	gopkg.in/yaml.v2 v2.2.4 // indirect
	k8s.io/client-go v11.0.0+incompatible

)

replace github.com/docker/distribution => github.com/docker/distribution v2.7.1+incompatible
