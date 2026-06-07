package repository

import (
	"database/sql"
	"memoria/internal/model"
)

type SessionRepo struct {
	DB *sql.DB
}

func (r *SessionRepo) Create(s model.Session) error {
	_, err := r.DB.Exec(`
		INSERT INTO sessions (id, user_id, title)
		VALUES ($1, $2, $3)
	`, s.ID, s.UserID, s.Title)

	return err
}

func (r *SessionRepo) ListByUser(userID string) ([]model.Session, error) {
	rows, err := r.DB.Query(`
		SELECT id, user_id, title, created_at
		FROM sessions
		WHERE user_id = $1
	`, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var sessions []model.Session

	for rows.Next() {
		var s model.Session
		rows.Scan(&s.ID, &s.UserID, &s.Title, &s.CreatedAt)
		sessions = append(sessions, s)
	}

	return sessions, nil
}
