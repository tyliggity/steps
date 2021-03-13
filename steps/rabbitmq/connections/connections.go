package main

import (
	"github.com/stackpulse/public-steps/common/step"
	"github.com/stackpulse/public-steps/rabbitmq/base"
)

type RMQConnections struct {
	*base.RabbitMQBase
}

func (r *RMQConnections) Init() error {
	rmqBase, err := base.NewRabbitMQBase(nil)
	if err != nil {
		return err
	}
	r.RabbitMQBase = rmqBase
	return nil
}

func (r *RMQConnections) Run() (int, []byte, error) {
	output, err := r.RunQuery("connections")
	if err != nil {
		return 1, nil, err
	}

	return 0, output, nil
}

func main() {
	step.Run(&RMQConnections{})
}
