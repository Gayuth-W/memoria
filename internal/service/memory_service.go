package service

import (
	"crypto/sha256"
	"encoding/hex"
	"memoria/internal/cache"
	"memoria/internal/model"
	"memoria/internal/repository"
	"memoria/internal/worker"

	"github.com/google/uuid"
)

type MemoryService struct {
	Repo   *repository.MemoryRepo
	Worker *worker.Worker
	Cache  *cache.RedisCache
}

func (s *MemoryService) Create(userID, sessionID, text string) error {

	hash := sha256.Sum256([]byte(text))
	embeddingHash := hex.EncodeToString(hash[:])

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
