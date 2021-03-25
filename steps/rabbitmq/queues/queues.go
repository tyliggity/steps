package main

import (
	"github.com/stackpulse/steps-sdk-go/step"
	"github.com/stackpulse/steps/rabbitmq/queues/base"
)

type RMQQueues struct {
	*base.RMQQueues
}

func (r *RMQQueues) Init() error {
	rmqQueue, err := base.NewRmqQueues()
	if err != nil {
		return err
	}
	r.RMQQueues = rmqQueue
	return nil
}

func (r *RMQQueues) Run() (int, []byte, error) {
	output, err := r.Get()
	if err != nil {
		return 1, nil, err
	}

	return 0, output, nil
}

func main() {
	step.Run(&RMQQueues{})
}
