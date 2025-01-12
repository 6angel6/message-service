package kafka

import (
	"context"
	"encoding/json"
	"github.com/segmentio/kafka-go"
	"log"
	"message-service/internal/model"
)

type Producer struct {
	w *kafka.Writer
}

func NewProducer(addr []string) *Producer {
	w := &kafka.Writer{
		Addr:                   kafka.TCP(addr...),
		Balancer:               &kafka.LeastBytes{},
		Topic:                  "unprocessed-messages",
		AllowAutoTopicCreation: true,
		RequiredAcks:           1,
		Logger:                 log.Default(),
		ErrorLogger:            log.Default(),
	}
	return &Producer{
		w: w,
	}
}

func (p *Producer) SendMessage(ctx context.Context, msgs ...model.Message) error {
	msgsToSend := make([]kafka.Message, 0, len(msgs))

	for _, msg := range msgs {
		JSON, err := json.Marshal(msg)
		if err != nil {
			return err
		}

		msgsToSend = append(msgsToSend, kafka.Message{
			Key:   []byte(msg.Id.String()),
			Value: JSON})
	}
	err := p.w.WriteMessages(ctx, msgsToSend...)
	if err != nil {
		return err
	}
	return nil
}
