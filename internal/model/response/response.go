package response

import (
	"github.com/gofrs/uuid"
	"time"
)

type MessageResponse struct {
	Id        uuid.UUID `json:"id"`
	Content   string    `json:"content"`
	Status    string    `json:"status"`
	CreatedAt time.Time `json:"created-at"`
}
