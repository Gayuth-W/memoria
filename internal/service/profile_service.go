package service

import (
	"github.com/google/uuid"
	"memoria/internal/repository"
)

type ProfileService struct {
	Repo *repository.ProfileRepo
}

func (s *ProfileService) GetProfile(userID uuid.UUID) ([]string, error) {
	return s.Repo.GetByUser(userID)
}

func (s *ProfileService) AddFact(userID uuid.UUID, fact string) error {
	return s.Repo.Add(userID, fact)
}

func (s *ProfileService) RemoveFact(userID uuid.UUID, fact string) error {
	return s.Repo.Remove(userID, fact)
}
