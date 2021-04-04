package main

import (
	"context"
	"encoding/json"
	"fmt"
	"reflect"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/ec2"
	envconf "github.com/caarlos0/env/v6"
	"github.com/stackpulse/steps-sdk-go/env"
	"github.com/stackpulse/steps-sdk-go/step"
)

const (
	MaxInstances = 1000
)

type Args struct {
	InstanceIDSJson env.JSONItemsArray `env:"INSTANCE_IDS,required" envDefault:""`
	Region          string             `env:"REGION,required" envDefault:""`
}

type DescribeInstances struct {
	args Args
}

func (s *DescribeInstances) Init() error {
	err := envconf.ParseWithFuncs(&s.args, map[reflect.Type]envconf.ParserFunc{
		reflect.TypeOf(env.JSONItemsArray{}): env.ParseJSONItemsArray,
	})
	if err != nil {
		return fmt.Errorf("init: %w", err)
	}

	if len(s.args.InstanceIDSJson.Items) > MaxInstances {
		return fmt.Errorf("INSTANCE_IDS parameter exceeded the max hosts allowed which is %d", MaxInstances)
	}

	return nil
}

type output struct {
	Ips []string `json:"ips"`
}

func (s *DescribeInstances) Run() (int, []byte, error) {
	ctx := context.Background()

	// Load the Shared AWS Configuration (~/.aws/config)
	cfg, err := config.LoadDefaultConfig(ctx, config.WithRegion(s.args.Region))
	if err != nil {
		return step.ExitCodeFailure, nil, fmt.Errorf("load default config: %w", err)
	}

	// Create an Amazon EC2 service client
	client := ec2.NewFromConfig(cfg)
	input := &ec2.DescribeInstancesInput{
		InstanceIds: s.args.InstanceIDSJson.Items,
	}

	result, err := client.DescribeInstances(ctx, input)
	if err != nil {
		return step.ExitCodeFailure, nil, fmt.Errorf("describe instances: %w", err)
	}

	var hostIPs output
	for _, reservation := range result.Reservations {
		for _, instance := range reservation.Instances {
			hostIPs.Ips = append(hostIPs.Ips, *instance.PrivateIpAddress)
		}
	}

	jsonOutput, err := json.Marshal(&hostIPs)
	if err != nil {
		return step.ExitCodeFailure, nil, fmt.Errorf("marshal output failed, error: %w", err)
	}

	return step.ExitCodeOK, jsonOutput, nil
}

func main() {
	step.Run(&DescribeInstances{})
}
