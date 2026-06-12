package model

import "time"

type Memory struct {
	ID              string    `json:"id"`
	UserID          string    `json:"user_id"`
	SessionID       string    `json:"session_id"`
	Text            string    `json:"text"`
	CreatedAt       time.Time `json:"created_at"`
	TSV             string    `json:"tsv"`
	ImportanceScore float64   `json:"importance_score"`
	EmbeddingHash   string    `json:"embedding_hash"`
}
