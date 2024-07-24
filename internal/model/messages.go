package model

import "time"

type Messages struct {
	UUID      int       `json:"-" db:"id"`
	Content   string    `json:"content" binding:"required" db:"content"`
	Status    string    `json:"status" db:"status"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
}
