package handler

import (
	"encoding/json"
	"net/http"

	"memoria/internal/middleware"
	"memoria/internal/service"
)

type MemoryHandler struct {
	Service *service.MemoryService
}

func (h *MemoryHandler) Create(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value(middleware.UserIDKey).(string)

	var req struct {
		SessionID string `json:"session_id"`
		Text      string `json:"text"`
	}

	json.NewDecoder(r.Body).Decode(&req)

	err := h.Service.Create(userID, req.SessionID, req.Text)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	w.WriteHeader(http.StatusCreated)
}
