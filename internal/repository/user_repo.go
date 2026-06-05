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
