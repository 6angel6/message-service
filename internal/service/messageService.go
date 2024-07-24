package service

import (
	"errors"
	"message-service/internal/model"
	"message-service/internal/repository"
	"time"
)

type MessageService struct {
	repo repository.Message
}

func NewMessageService(repo *repository.Repository) *MessageService {
	return &MessageService{repo: repo}
}

func (s *MessageService) ProcessMessage(content string) error {
	if content == "" {
		return errors.New("content is required")
	}
	msg := &model.Messages{
		Content:   content,
		Status:    "pending",
		CreatedAt: time.Now(),
	}
	return s.repo.SaveMessage(msg)
}
