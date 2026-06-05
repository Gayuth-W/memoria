package repository

import (
	"database/sql"
	"memoria/internal/model"
)

type MemoryRepo struct {
	DB *sql.DB
}

func (r *MemoryRepo) Create(m model.Memory) error {
	_, err := r.DB.Exec(`
		INSERT INTO memories (id, user_id, session_id, text)
		VALUES ($1, $2, $3, $4)
	`, m.ID, m.UserID, m.SessionID, m.Text)

	return err
}

func (r *MemoryRepo) ListByUser(userID string) ([]model.Memory, error) {
	rows, err := r.DB.Query(`
		SELECT id, user_id, session_id, text, created_at
		FROM memories
		WHERE user_id = $1
	`, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var res []model.Memory

	for rows.Next() {
		var m model.Memory
		rows.Scan(&m.ID, &m.UserID, &m.SessionID, &m.Text, &m.CreatedAt)
		res = append(res, m)
	}

	return res, nil
}
