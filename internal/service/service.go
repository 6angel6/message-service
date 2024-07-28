package service

import (
	"context"
	"message-service/internal/kafka"
	"message-service/internal/model"
	"message-service/internal/repository"
)

type Message interface {
	CreateMessage(content string) error
	GetAllMessages() ([]model.Message, error)
	ProcessMsg(ctx context.Context, msg model.Message) error
	SendMsgToKafka(ctx context.Context) error
}

type Producer interface {
	SendMessage(ctx context.Context, msgs ...model.Message) error
}

type Service struct {
	Message
	Producer
}

func NewService(repo *repository.Repository, producer *kafka.Producer) *Service {
	return &Service{
		Message: NewMessageService(repo.Message, producer),
	}
}
