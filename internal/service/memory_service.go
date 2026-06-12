package service

import (
	"memoria/internal/model"
	"memoria/internal/repository"
	"memoria/internal/worker"

	"github.com/google/uuid"
)

type MemoryService struct {
	Repo   *repository.MemoryRepo
	Worker *worker.Worker
}

func (s *MemoryService) Create(userID, sessionID, text, embeddingHash string) error {

	m := model.Memory{
		ID:            uuid.New().String(),
		UserID:        userID,
		SessionID:     sessionID,
		Text:          text,
		EmbeddingHash: embeddingHash,
	}

	// 1. save to postgres
	err := s.Repo.Create(m)
	if err != nil {
		return err
	}

	// 2. async embedding job
	s.Worker.Enqueue(worker.Job{
		MemoryID:  m.ID,
		UserID:    userID,
		SessionID: sessionID,
		Text:      text,
	})

	return nil
}
