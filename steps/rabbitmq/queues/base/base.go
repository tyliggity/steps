package base

import (
	"fmt"
	"net/url"
	"strings"

	"github.com/stackpulse/steps/rabbitmq/base"
)

type Args struct {
	*base.Args
	Vhost     string `env:"VHOST"`
	QueueName string `env:"QUEUE_NAME"`
}

type RMQQueues struct {
	*base.RabbitMQBase
	args *Args
}

func NewRmqQueues() (*RMQQueues, error) {
	args := &Args{Args: &base.Args{}}
	rmqBase, err := base.NewRabbitMQBase(args)
	if err != nil {
		return nil, err
	}

	return &RMQQueues{RabbitMQBase: rmqBase, args: args}, nil
}

func (r *RMQQueues) Get() ([]byte, error) {
	var sb strings.Builder
	sb.WriteString("queues")
	if r.args.Vhost != "" {
		sb.WriteString(fmt.Sprintf("/%s", url.QueryEscape(r.args.Vhost)))
		if r.args.QueueName != "" {
			sb.WriteString(fmt.Sprintf("/%s", url.QueryEscape(r.args.QueueName)))
		}
	}

	output, err := r.RunQuery(sb.String())
	if err != nil {
		return nil, err
	}

	return output, nil
}
