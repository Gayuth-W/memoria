package service

import (
	"memoria/internal/model"
	"memoria/internal/repository"

	"github.com/google/uuid"
)

type SessionService struct {
	Repo       *repository.SessionRepo
	MemoryRepo *repository.MemoryRepo
}

func (s *SessionService) Create(userID uuid.UUID, title string) error {
	return s.Repo.Create(model.Session{
		ID:     uuid.New().String(),
		UserID: userID,
		Title:  title,
	})
}

func (s *SessionService) ListByUser(userID string) ([]model.Session, error) {
	return s.Repo.ListByUser(userID)
}

func (s *SessionService) GetByID(id, userID string) (*model.Session, error) {
	return s.Repo.GetByID(id, userID)
}

func (s *SessionService) GetMemories(sessionID, userID string) ([]model.Memory, error) {
	// First check if the session belongs to the user
	session, err := s.Repo.GetByID(sessionID, userID)
	if err != nil {
		return nil, err
	}
	if session == nil {
		return nil, repository.ErrNotFound
	}

	if s.MemoryRepo != nil {
		return s.MemoryRepo.ListBySession(sessionID, userID)
	}
	return nil, nil
}
