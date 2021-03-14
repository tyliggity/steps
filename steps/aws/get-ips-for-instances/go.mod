module github.com/stackpulse/public-steps/aws/command

go 1.14

require (
	github.com/aws/aws-sdk-go v1.36.32 // indirect
	github.com/aws/aws-sdk-go-v2 v1.0.0 // indirect
	github.com/aws/aws-sdk-go-v2/config v1.0.0
	github.com/aws/aws-sdk-go-v2/service/ec2 v1.0.0
	github.com/aws/aws-sdk-go-v2/service/s3 v1.0.0 // indirect
	github.com/caarlos0/env/v6 v6.5.0
	github.com/stackpulse/public-steps/common v0.0.0
)

replace github.com/stackpulse/public-steps/common v0.0.0 => ../../common
