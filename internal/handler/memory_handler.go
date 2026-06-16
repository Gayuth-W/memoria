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
		SessionID string `json:"session_id"`
		Text      string `json:"text"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	if req.Text == "" {
		http.Error(w, "text is required", http.StatusBadRequest)
		return
	}

	if req.SessionID == "" {
		http.Error(w, "session_id is required", http.StatusBadRequest)
		return
	}

	err := h.Service.Create(userID, req.SessionID, req.Text)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
}
