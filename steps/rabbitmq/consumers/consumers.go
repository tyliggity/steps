package main

import (
	"github.com/stackpulse/public-steps/rabbitmq/base"
	"github.com/stackpulse/steps-sdk-go/step"
)

type RMQConsumers struct {
	*base.RabbitMQBase
}

func (r *RMQConsumers) Init() error {
	rmqBase, err := base.NewRabbitMQBase(nil)
	if err != nil {
		return err
	}
	r.RabbitMQBase = rmqBase
	return nil
}

func (r *RMQConsumers) Run() (int, []byte, error) {
	output, err := r.RunQuery("consumers")
	if err != nil {
		return 1, nil, err
	}

	return 0, output, nil
}

func main() {
	step.Run(&RMQConsumers{})
}
