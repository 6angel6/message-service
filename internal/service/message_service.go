package service

import (
	"context"
	"errors"
	"github.com/gofrs/uuid"
	"message-service/internal/config"
	"message-service/internal/kafka"
	"message-service/internal/model"
	"message-service/internal/repository"
	"time"
)

type MessageService struct {
	repo     repository.Message
	producer *kafka.Producer
}

func NewMessageService(repo repository.Message, producer *kafka.Producer) *MessageService {
	return &MessageService{repo: repo,
		producer: producer}
}

func (s *MessageService) CreateMessage(content string) error {
	if content == "" {
		return errors.New("content is required")
	}
	msg := &model.Message{
		Id:        uuid.Must(uuid.NewV4()),
		Content:   content,
		Status:    config.PENDING,
		CreatedAt: time.Now(),
	}

	return s.repo.CreateMessage(msg)
}

func (s *MessageService) SendMsgToKafka(ctx context.Context) error {
	msgs, err := s.repo.UnprocessedMsgs()
	if err != nil {
		return err
	}

	err = s.repo.UpdateMessageStatus(config.PROCESSING, msgs...)

	time.Sleep(10 * time.Second)

	err = s.producer.SendMessage(ctx, msgs...)
	if err != nil {
		return err
	}

	return nil
}

func (s *MessageService) ProcessMsg(_ context.Context, msg model.Message) error {
	return s.repo.UpdateMessageStatus(config.PROCESSED, msg)
}

func (s *MessageService) GetAllMessages() ([]model.Message, error) {
	return s.repo.GetAllMessages()
}
