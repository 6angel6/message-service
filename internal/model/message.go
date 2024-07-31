package model

import (
	"github.com/gofrs/uuid"
	"time"
)

type Message struct {
	Id        uuid.UUID `json:"id" db:"id"`
	Content   string    `json:"content" binding:"required" db:"content"`
	Status    string    `json:"status" db:"status"`
	CreatedAt time.Time `json:"-" db:"created_at"`
}
