package handler

import (
	"encoding/json"
	"net/http"

	"memoria/internal/middleware"
	"memoria/internal/service"

	"github.com/google/uuid"
)

type MemoryHandler struct {
	Service *service.MemoryService
}

func (h *MemoryHandler) Create(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value(middleware.UserIDKey).(uuid.UUID).String()

	var req struct {
		SessionID     string `json:"session_id"`
		Text          string `json:"text"`
		EmbeddingHash string `json:"embedding_hash"`
	}

	json.NewDecoder(r.Body).Decode(&req)

	err := h.Service.Create(userID, req.SessionID, req.Text, req.EmbeddingHash)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	w.WriteHeader(http.StatusCreated)
}
