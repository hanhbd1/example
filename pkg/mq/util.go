package mq

import (
	"encoding/json"
	"fmt"

	"github.com/ThreeDotsLabs/watermill/message"
	"github.com/google/uuid"
)

func Bind[B interface{}](msg *message.Message) (B, error) {
	var body B
	if err := json.Unmarshal(msg.Payload, &body); err != nil {
		return body, fmt.Errorf("mqutil: failed bind %v", err)
	}
	return body, nil
}

func Publish(publisher message.Publisher, topic string, data []byte) error {
	msg := message.NewMessage(uuid.NewString(), data)
	return publisher.Publish(topic, msg)
}
