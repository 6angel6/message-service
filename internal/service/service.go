package service

import (
	"context"
	"message-service/internal/kafka"
	"message-service/internal/model"
	"message-service/internal/model/response"
	"message-service/internal/repository"
)

type Message interface {
	CreateMessage(content string) error
	ProcessMsg(ctx context.Context, msg model.Message) error
	SendMsgToKafka(ctx context.Context) error
	GetStats() ([]response.MessageResponse, error)
	ConvertToMessageResponse(msg model.Message) response.MessageResponse
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
