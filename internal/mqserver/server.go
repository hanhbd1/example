package mqserver

import (
	"context"

	"example/configreader"
	"example/pkg/log"
	"example/pkg/mq"
	"example/pkg/mq/publisher"
	"example/pkg/mq/subscriber"

	"github.com/ThreeDotsLabs/watermill/message"
)

// Instance represents an instance of the server
type Instance struct {
	ctx        context.Context
	route      *message.Router
	subscriber message.Subscriber
	publisher  message.Publisher
}

// NewInstance returns a new instance of our server
func NewInstance() *Instance {
	return &Instance{}
}

// Start starts the server
func (i *Instance) Start(ctx context.Context) {
	i.ctx = ctx
	// Create a new router
	var err error
	i.route, err = message.NewRouter(message.RouterConfig{}, mq.NewLogger())
	if err != nil {
		log.Fatalw("failed to init consumer engine", "error", err.Error())
	}
	i.subscriber, err = subscriber.New(configreader.Config.Queue.Name)
	i.publisher, err = publisher.New()
	// example for publisher
	// bb, err := &dto.SimpleMessage{Message: "hello"}.ToBytes()
	// mq.Publish(i.publisher, simpleTopic, bb)
	InitRoute(i.route, i.subscriber)
	err = i.route.Run(i.ctx)
	if err != nil {
		log.Errorw("consumer stopped unexpected", "error", err)
		i.Shutdown(ctx)
	}
}

// Shutdown stops the server
func (i *Instance) Shutdown(ctx context.Context) {
	i.subscriber.Close()
	i.route.Close()
}
