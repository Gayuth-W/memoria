package service

import (
	"memoria/internal/model"
	"memoria/internal/repository"

	"github.com/google/uuid"
)

type SessionService struct {
	Repo *repository.SessionRepo
}

func (s *SessionService) Create(userID uuid.UUID, title string) error {
	return s.Repo.Create(model.Session{
		ID:     uuid.New().String(),
		UserID: userID,
		Title:  title,
	})
}
