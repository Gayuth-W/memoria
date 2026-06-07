package repository

import (
	"database/sql"
	"memoria/internal/model"
)

type UserRepo struct {
	DB *sql.DB
}

func (r *UserRepo) GetByAPIKey(key string) (*model.User, error) {
	u := &model.User{}

	err := r.DB.QueryRow(`
		SELECT id, api_key, created_at
		FROM users
		WHERE api_key = $1
	`, key).Scan(&u.ID, &u.APIKey, &u.CreatedAt)

	return u, err
}

func (r *UserRepo) Create(s model.User) error {
	_, err := r.DB.Exec(`
		insert into users(id, api_key, created_at)
		values ($1, $2, $3)
	`, s.ID, s.APIKey, s.CreatedAt)

	return err
}
