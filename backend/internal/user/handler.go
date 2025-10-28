package user

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

func (h *Handler) GetAllUsers(w http.ResponseWriter, r *http.Request) {
	users, err := h.service.GetAllUsers()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(users)
}

func (h *Handler) CreateUser(w http.ResponseWriter, r *http.Request) {
	// Implementation needed
	w.WriteHeader(http.StatusCreated)
}

func (h *Handler) UpdateUser(w http.ResponseWriter, r *http.Request) {
	// Implementation needed
	w.WriteHeader(http.StatusOK)
}

func (h *Handler) ChangeUserPassword(w http.ResponseWriter, r *http.Request) {
	// Implementation needed
	w.WriteHeader(http.StatusOK)
}

func (h *Handler) DeleteUser(w http.ResponseWriter, r *http.Request) {
	// Implementation needed
	w.WriteHeader(http.StatusOK)
}

func (h *Handler) GetAllRoles(w http.ResponseWriter, r *http.Request) {
	roles, err := h.service.GetAllRoles()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(roles)
}

func (h *Handler) CreateRole(w http.ResponseWriter, r *http.Request) {
	// Implementation needed
	w.WriteHeader(http.StatusCreated)
}

func (h *Handler) UpdateRole(w http.ResponseWriter, r *http.Request) {
	// Implementation needed
	w.WriteHeader(http.StatusOK)
}

func (h *Handler) DeleteRole(w http.ResponseWriter, r *http.Request) {
	// Implementation needed
	w.WriteHeader(http.StatusOK)
}
