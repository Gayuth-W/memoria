package model

import (
	"time"

	"github.com/google/uuid"
)

type PinnedFact struct {
	ID        uuid.UUID `json:"id"`
	UserID    uuid.UUID `json:"user_id"`
	Fact      string    `json:"fact"`
	CreatedAt time.Time `json:"created_at"`
}
