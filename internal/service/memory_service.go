package service

import (
	"context"
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

	if s.Cache != nil {
		// A simple way to invalidate search cache is to use Redis keys pattern
		// For production, we'd use SETs or tags, but let's just clear matching keys.
		// Redis SCAN would be better, but for simplicity we can use Keys if not huge, or just a known prefix
		ctx := context.Background()
		keys, err := s.Cache.Client.Keys(ctx, "search:*").Result()
		if err == nil && len(keys) > 0 {
			s.Cache.Client.Del(ctx, keys...)
		}
	}
	return nil
}

func (s *MemoryService) ListByUser(userID string) ([]model.Memory, error) {
	return s.Repo.ListByUser(userID)
}

func (s *MemoryService) Delete(id string) error {
	return s.Repo.Delete(id)
}

func (s *MemoryService) GetByID(id, userID string) (*model.Memory, error) {
	return s.Repo.GetByID(id, userID)
}
