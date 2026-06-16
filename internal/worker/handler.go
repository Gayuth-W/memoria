package worker

import (
	"memoria/internal/embedding"
	vector "memoria/internal/qdrant"
	"sync"
)

type Handler struct {
	Embedder embedding.Embedder
	Vector   *vector.VectorStore
	cache    sync.Map
}

func (h *Handler) Handle(job Job) error {
	var vec []float32
	if cached, ok := h.cache.Load(job.Text); ok {
		vec = cached.([]float32)
	} else {
		var err error
		vec, err = h.Embedder.Embed(job.Text)
		if err != nil {
			return err
		}
		h.cache.Store(job.Text, vec)
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
