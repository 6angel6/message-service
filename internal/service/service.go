package service

import (
	"message-service/internal/repository"
)

type Message interface {
	ProcessMessage(content string) error
}

type Service struct {
	Message
}

func NewService(repo *repository.Repository) *Service {
	return &Service{
		Message: NewMessageService(repo),
	}
}
