package handler

import (
	"encoding/json"
	"net/http"

	"memoria/internal/middleware"
	"memoria/internal/search"

	"github.com/google/uuid"
)

type SearchHandler struct {
	Service *search.Service
}

func (h *SearchHandler) Search(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value(middleware.UserIDKey).(uuid.UUID).String()

	var req struct {
		SessionID string `json:"session_id"`
		Query     string `json:"query"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid body", http.StatusBadRequest)
		return
	}

	results, trace, err := h.Service.Search(userID, req.SessionID, req.Query)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	resp := map[string]any{"results": results}
	if r.URL.Query().Get("debug") == "true" {
		resp["trace"] = trace
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}
