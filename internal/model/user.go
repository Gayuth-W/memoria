package model

import (
	"time"

	"github.com/google/uuid"
)

type User struct {
	ID        uuid.UUID `json:"id"`
	APIKey    string    `json:"api_key"`
	CreatedAt time.Time `json:"created_at"`
}
