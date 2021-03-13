module github.com/stackpulse/steps/flux/lock

go 1.15

require (
	github.com/caarlos0/env/v6 v6.5.0
	github.com/docker/distribution v2.7.1+incompatible // indirect
	github.com/fluxcd/flux v1.17.0
	github.com/fluxcd/flux/pkg/install v0.0.0-20210210162146-6a501128a692 // indirect
	github.com/opencontainers/go-digest v1.0.0 // indirect
	github.com/opencontainers/image-spec v1.0.1 // indirect
	github.com/stackpulse/steps/common v0.0.0
	k8s.io/client-go v11.0.0+incompatible

)

replace (
	github.com/docker/distribution => github.com/docker/distribution v2.7.1+incompatible
	github.com/stackpulse/steps/common v0.0.0 => ../../common
)
