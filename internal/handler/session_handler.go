package handler

import (
	"encoding/json"
	"net/http"

	"memoria/internal/middleware"
	"memoria/internal/service"

	"github.com/google/uuid"
)

type SessionHandler struct {
	Service *service.SessionService
}

func (h *SessionHandler) Create(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value(middleware.UserIDKey).(uuid.UUID)

	var req struct {
		Title string `json:"title"`
	}

	json.NewDecoder(r.Body).Decode(&req)

	err := h.Service.Create(userID, req.Title)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

func (h *SessionHandler) ListByUser(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value(middleware.UserIDKey).(uuid.UUID).String()

	sessions, err := h.Service.ListByUser(userID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(sessions)
}
