package mqserver

import (
	"example/internal/mqhandler"

	"github.com/ThreeDotsLabs/watermill/message"
)

func InitRoute(r *message.Router, subscriber message.Subscriber) {
	mqHandler := mqhandler.New()
	r.AddNoPublisherHandler("SimpleMessageHandler", simpleTopic, subscriber, mqHandler.SimpleMessageHandler)
}
