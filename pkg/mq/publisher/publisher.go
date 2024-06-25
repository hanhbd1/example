package publisher

import (
	"fmt"

	"example/configreader"
	"example/pkg/mq"

	"github.com/ThreeDotsLabs/watermill"
	"github.com/ThreeDotsLabs/watermill-amqp/pkg/amqp"
	"github.com/ThreeDotsLabs/watermill/message"
)

func New() (message.Publisher, error) {
	logger := mq.NewLogger()
	provider := configreader.Config.Queue.Provider

	var p message.Publisher
	var err error
	switch provider {
	case "ampq":
		p, err = newAmqpPublisher(logger)
	default:
		return nil, fmt.Errorf("not support provider %s", provider)
	}
	if err != nil {
		return nil, err
	}
	return p, nil
}

func newAmqpPublisher(logger watermill.LoggerAdapter) (message.Publisher, error) {
	config := amqp.NewDurablePubSubConfig(configreader.Config.Queue.Amqp, nil)
	return amqp.NewPublisher(config, logger)
}
