package mqhandler

import (
	"example/internal/dto"
	"example/pkg/log"
	"example/pkg/mq"

	"github.com/ThreeDotsLabs/watermill/message"
)

func (h *Handler) SimpleMessageHandler(msg *message.Message) error {
	payload, err := mq.Bind[dto.SimpleMessage](msg)
	if err != nil {
		msg.Ack()
		return err
	}
	log.Infow("SimpleMessageHandler", "message", payload.Message)
	return nil
}
