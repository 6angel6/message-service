package kafka

import (
	"context"
	"encoding/json"
	"github.com/segmentio/kafka-go"
	"log"
	"message-service/internal/model"
)

type Consumer struct {
	r *kafka.Reader
}

func NewConsumer(brokers []string) *Consumer {
	r := kafka.NewReader(kafka.ReaderConfig{
		Brokers:     brokers,
		Topic:       "unprocessed-messages",
		GroupID:     "msg-processor",
		MaxBytes:    10e3,
		Logger:      log.Default(),
		ErrorLogger: log.Default(),
	})
	return &Consumer{r: r}
}

func (c *Consumer) ConsumeMessage(ctx context.Context, fn func(context.Context, model.Message) error) {
	for {
		select {
		case <-ctx.Done():
			return
		default:
		}
		m, err := c.r.ReadMessage(ctx)
		if err != nil {
			log.Println("Error reading message:", err)
			continue
		}

		var msg model.Message
		err = json.Unmarshal(m.Value, &msg)
		if err != nil {
			log.Println("Error unmarshalling message:", err)
			continue
		}

		err = fn(ctx, msg)
		if err != nil {
			log.Println("Error processing message:", err)
			continue
		}

	}
}
