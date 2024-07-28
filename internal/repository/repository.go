package repository

import (
	"database/sql"
	"message-service/internal/model"
)

type Message interface {
	CreateMessage(msg *model.Message) error
	GetAllMessages() ([]model.Message, error)
	UpdateMessageStatus(status string, msgs ...model.Message) error
	UnprocessedMsgs() ([]model.Message, error)
}

type Repository struct {
	Message
}

func NewRepository(db *sql.DB) *Repository {
	return &Repository{
		Message: NewMessageRepo(db),
	}
}
