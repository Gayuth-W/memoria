package service

import (
	"memoria/internal/model"
	"memoria/internal/repository"

	"github.com/google/uuid"
)

type MemoryService struct {
	Repo *repository.MemoryRepo
}

func (s *MemoryService) Create(userID, sessionID, text string) error {
	return s.Repo.Create(model.Memory{
		ID:        uuid.New().String(),
		UserID:    userID,
		SessionID: sessionID,
		Text:      text,
	})
}
