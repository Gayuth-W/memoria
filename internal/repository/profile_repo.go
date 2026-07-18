package repository

import (
	"database/sql"
	"time"

	"github.com/google/uuid"
)

type ProfileRepo struct {
	DB *sql.DB
}

func (r *ProfileRepo) GetByUser(userID uuid.UUID) ([]string, error) {
	rows, err := r.DB.Query(`
		SELECT fact FROM pinned_facts
		WHERE user_id = $1
		ORDER BY created_at ASC
	`, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var facts []string
	for rows.Next() {
		var f string
		if err := rows.Scan(&f); err != nil {
			return nil, err
		}
		facts = append(facts, f)
	}

	return facts, nil
}

func (r *ProfileRepo) Add(userID uuid.UUID, fact string) error {
	id := uuid.New()
	_, err := r.DB.Exec(`
		INSERT INTO pinned_facts(id, user_id, fact, created_at)
		VALUES ($1, $2, $3, $4)
		ON CONFLICT (user_id, fact) DO NOTHING
	`, id, userID, fact, time.Now())
	return err
}

func (r *ProfileRepo) Remove(userID uuid.UUID, fact string) error {
	_, err := r.DB.Exec(`
		DELETE FROM pinned_facts
		WHERE user_id = $1 AND fact = $2
	`, userID, fact)
	return err
}
