package auth

import (
	"encoding/json"
	"net/http"

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
	// For now, we'll just call the service
	token, err := h.service.Login("test", "test")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(map[string]string{"token": token})
}

// @Summary Logout a user
// @Description Logout a user
// @Tags auth
// @Success 200
// @Router /auth/logout [post]
func (h *Handler) Logout(w http.ResponseWriter, r *http.Request) {
	err := h.service.Logout("test")
	if err != nil {
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
	user, err := h.service.GetMe("test")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(user)
}

type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}
