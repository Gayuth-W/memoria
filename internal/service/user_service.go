package service

import (
	"memoria/internal/model"
	"memoria/internal/repository"
	"time"

	"github.com/google/uuid"
)

type UserService struct {
	Repo *repository.UserRepo
}

func (s *UserService) Create(APIKey string) error {
	return s.Repo.Create(model.User{
		ID:        uuid.New(),
		APIKey:    APIKey,
		CreatedAt: time.Now(),
	})
}
