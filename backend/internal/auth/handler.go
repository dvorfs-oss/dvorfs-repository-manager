package auth

import (
	"encoding/json"
	"errors"
	"net/http"
	"strings"

	_ "dvorfs-repository-manager/internal/user"
)

type Handler struct {
	service Service
}

func NewHandler(service Service) *Handler {
	return &Handler{service}
}

// @Summary Login a user
// @Description Login a user with username and password
// @Tags auth
// @Accept  json
// @Produce  json
// @Param   credentials body LoginRequest true "Credentials"
// @Success 200 {object} map[string]string
// @Router /auth/login [post]
func (h *Handler) Login(w http.ResponseWriter, r *http.Request) {
	var req LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid credentials payload", http.StatusBadRequest)
		return
	}
	token, err := h.service.Login(req.Username, req.Password)
	if err != nil {
		if errors.Is(err, ErrInvalidCredentials) {
			http.Error(w, err.Error(), http.StatusUnauthorized)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"token": token})
}

// @Summary Logout a user
// @Description Logout a user
// @Tags auth
// @Success 200
// @Router /auth/logout [post]
func (h *Handler) Logout(w http.ResponseWriter, r *http.Request) {
	token := strings.TrimSpace(strings.TrimPrefix(r.Header.Get("Authorization"), "Bearer "))
	err := h.service.Logout(token)
	if err != nil {
		if errors.Is(err, ErrInvalidToken) {
			http.Error(w, err.Error(), http.StatusUnauthorized)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}

// @Summary Get current user
// @Description Get the current authenticated user
// @Tags auth
// @Produce  json
// @Success 200 {object} user.User
// @Router /auth/me [get]
func (h *Handler) GetMe(w http.ResponseWriter, r *http.Request) {
	token := strings.TrimSpace(strings.TrimPrefix(r.Header.Get("Authorization"), "Bearer "))
	user, err := h.service.GetMe(token)
	if err != nil {
		if errors.Is(err, ErrInvalidToken) {
			http.Error(w, err.Error(), http.StatusUnauthorized)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(user)
}

type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}
