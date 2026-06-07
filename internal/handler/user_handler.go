package handler

import (
	"encoding/json"
	"net/http"

	"memoria/internal/service"
)

type UserHandler struct {
	Service *service.UserService
}

func (h *UserHandler) Create(w http.ResponseWriter, r *http.Request) {

	var req struct {
		APIKey string `json:"api_key"`
	}

	json.NewDecoder(r.Body).Decode(&req)

	err := h.Service.Create(req.APIKey)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	w.WriteHeader(http.StatusCreated)
}
