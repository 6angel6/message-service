package repository

import (
	"database/sql"
	"message-service/internal/model"
)

type Message interface {
	SaveMessage(msg *model.Messages) error
}

type Repository struct {
	Message
}

func NewRepository(db *sql.DB) *Repository {
	return &Repository{
		Message: NewMessageRepo(db),
	}
}
