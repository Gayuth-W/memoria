package repository

import (
	"database/sql"

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
