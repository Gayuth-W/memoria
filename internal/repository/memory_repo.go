package repository

import (
	"database/sql"
	"memoria/internal/model"

	"github.com/lib/pq"
)

type MemoryRepo struct {
	DB *sql.DB
}

func (r *MemoryRepo) Create(m model.Memory) error {
	_, err := r.DB.Exec(`
		INSERT INTO memories (
			id,
			user_id,
			session_id,
			text,
			embedding_hash
		)
		VALUES ($1,$2,$3,$4,$5)
	`,
		m.ID,
		m.UserID,
		m.SessionID,
		m.Text,
		m.EmbeddingHash,
	)

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

func (r *MemoryRepo) KeywordSearch(userID, query string) ([]string, error) {
	rows, err := r.DB.Query(`
		SELECT id
		FROM memories
		WHERE user_id = $1 AND text ILIKE $2
	`, userID, "%"+query+"%")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var ids []string
	for rows.Next() {
		var id string
		if err := rows.Scan(&id); err != nil {
			return nil, err
		}
		ids = append(ids, id)
	}

	return ids, nil
}

func (r *MemoryRepo) GetByIDs(
	ids []string,
) ([]model.Memory, error) {

	if len(ids) == 0 {
		return []model.Memory{}, nil
	}

	query := `
        SELECT
            id,
            user_id,
            session_id,
            text,
            created_at,
            importance_score
        FROM memories
        WHERE id = ANY($1)
    `

	rows, err := r.DB.Query(query, pq.Array(ids))
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var memories []model.Memory

	for rows.Next() {
		var m model.Memory

		err := rows.Scan(
			&m.ID,
			&m.UserID,
			&m.SessionID,
			&m.Text,
			&m.CreatedAt,
			&m.ImportanceScore,
		)
		if err != nil {
			return nil, err
		}

		memories = append(memories, m)
	}

	return memories, nil
}
