package repository

import (
	"database/sql"
	"github.com/gofrs/uuid"
	"github.com/lib/pq"
	_ "github.com/lib/pq"
	"log"
	"message-service/internal/config"
	"message-service/internal/model"
)

type MessageRepo struct {
	db *sql.DB
}

func NewMessageRepo(db *sql.DB) *MessageRepo {
	return &MessageRepo{db: db}
}

func (r *MessageRepo) CreateMessage(msg *model.Message) error {
	query := `INSERT INTO messages (content, status, created_at) VALUES ($1, $2, $3) RETURNING id`
	err := r.db.QueryRow(query, msg.Content, msg.Status, msg.CreatedAt).Scan(&msg.Id)
	if err != nil {
		log.Printf("Error creating message: %v", err)
		return err
	}
	return nil
}

func (r *MessageRepo) UpdateMessageStatus(status string, msgs ...model.Message) error {
	query := `UPDATE messages SET status = $1 WHERE id = ANY($2)`

	// Создаем массив UUID
	ids := make([]uuid.UUID, 0, len(msgs))
	for _, msg := range msgs {
		ids = append(ids, msg.Id)
	}

	// Преобразуем массив UUID в формат, который PostgreSQL может интерпретировать
	idsArray := pq.Array(ids)
	log.Printf("Updating status to %s for IDs: %v", status, idsArray)
	// Выполняем запрос
	_, err := r.db.Exec(query, status, idsArray)
	if err != nil {
		log.Printf("Error updating message status: %v", err)
		return err
	}

	return nil
}

func (r *MessageRepo) UnprocessedMsgs() ([]model.Message, error) {
	var messages []model.Message
	query := `SELECT id, content, status, created_at FROM messages WHERE status = $1`
	rows, err := r.db.Query(query, config.PENDING)
	if err != nil {
		log.Printf("Error fetching unprocessed messages: %v", err)
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var msg model.Message
		err := rows.Scan(&msg.Id, &msg.Content, &msg.Status, &msg.CreatedAt)
		if err != nil {
			log.Printf("Error scanning message: %v", err)
			return nil, err
		}
		messages = append(messages, msg)
	}

	if err = rows.Err(); err != nil {
		log.Printf("Error iterating over rows: %v", err)
		return nil, err
	}

	return messages, nil
}

func (r *MessageRepo) GetAllMessages() ([]model.Message, error) {
	query := `SELECT id, content, status, created_at FROM messages`
	rows, err := r.db.Query(query)
	if err != nil {
		return nil, err
	}

	defer rows.Close()
	var msgs []model.Message
	for rows.Next() {
		var msg model.Message
		if err := rows.Scan(&msg.Id, &msg.Content, &msg.Status, &msg.CreatedAt); err != nil {
			log.Printf("Error scanning row: %v", err)
			return nil, err
		}
		msgs = append(msgs, msg)
	}
	return msgs, nil
}
