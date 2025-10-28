package cleanup

import (
	"encoding/json"
	"net/http"
)

type Handler struct {
	service Service
}

func NewHandler(service Service) *Handler {
	return &Handler{service}
}

func (h *Handler) GetAllCleanupPolicies(w http.ResponseWriter, r *http.Request) {
	policies, err := h.service.GetAllCleanupPolicies()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(policies)
}

func (h *Handler) CreateCleanupPolicy(w http.ResponseWriter, r *http.Request) {
	// Implementation needed
	w.WriteHeader(http.StatusCreated)
}

func (h *Handler) UpdateCleanupPolicy(w http.ResponseWriter, r *http.Request) {
	// Implementation needed
	w.WriteHeader(http.StatusOK)
}

func (h *Handler) DeleteCleanupPolicy(w http.ResponseWriter, r *http.Request) {
	// Implementation needed
	w.WriteHeader(http.StatusOK)
}
