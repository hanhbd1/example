package subscriber

import (
	"errors"

	"example/configreader"
	"example/pkg/mq"

	"github.com/ThreeDotsLabs/watermill-amqp/pkg/amqp"
	"github.com/ThreeDotsLabs/watermill/message"
)

func New(name string) (message.Subscriber, error) {
	logger := mq.NewLogger()
	switch configreader.Config.Queue.Provider {
	case "amqp":
		config := amqp.NewDurablePubSubConfig(
			configreader.Config.Queue.Amqp,
			func(topic string) string {
				return topic + ":" + name
			},
		)

		return amqp.NewSubscriber(config, logger)
	}
	return nil, errors.New("provider not supported")
}
