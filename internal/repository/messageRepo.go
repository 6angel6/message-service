package repository

import (
	"database/sql"
	_ "github.com/lib/pq"
	"message-service/internal/model"
)

type MessageRepo struct {
	db *sql.DB
}

func NewMessageRepo(db *sql.DB) *MessageRepo {
	return &MessageRepo{db: db}
}

func (r *MessageRepo) SaveMessage(msg *model.Messages) error {
	query := `INSERT INTO messages (content, status, created_at) VALUES ($1, $2, $3) RETURNING id`
	return r.db.QueryRow(query, msg.Content, msg.Status, msg.CreatedAt).Scan(&msg.UUID)
}
