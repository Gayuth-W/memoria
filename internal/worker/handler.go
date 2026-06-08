package worker

import (
	"memoria/internal/embedding"
	"memoria/internal/vector"
)

type Handler struct {
	Embedder embedding.Embedder
	Vector   *vector.VectorStore
}

func (h *Handler) Handle(job Job) error {

	vec, err := h.Embedder.Embed(job.Text)
	if err != nil {
		return err
	}

	return h.Vector.Upsert(
		job.MemoryID,
		vec,
		map[string]any{
			"user_id":    job.UserID,
			"session_id": job.SessionID,
		},
	)
}
