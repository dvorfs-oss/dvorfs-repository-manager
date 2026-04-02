package auth

import (
	"net/http"
	"strings"

	"dvorfs-repository-manager/internal/user"
	"dvorfs-repository-manager/pkg/httpx"
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
	var request LoginRequest
	if err := httpx.DecodeJSON(r, &request); err != nil {
		httpx.WriteError(w, http.StatusBadRequest, err.Error())
		return
	}

	token, err := h.service.Login(request.Username, request.Password)
	if err != nil {
		httpx.WriteError(w, http.StatusUnauthorized, err.Error())
		return
	}

	account, err := h.service.GetMe(token)
	if err != nil {
		httpx.WriteError(w, http.StatusInternalServerError, err.Error())
		return
	}

	httpx.WriteJSON(w, http.StatusOK, map[string]any{
		"token": token,
		"user":  user.SanitizeUser(*account),
	})
}

// @Summary Logout a user
// @Description Logout a user
// @Tags auth
// @Success 200
// @Router /auth/logout [post]
func (h *Handler) Logout(w http.ResponseWriter, r *http.Request) {
	if err := h.service.Logout(extractBearerToken(r.Header.Get("Authorization"))); err != nil {
		httpx.WriteError(w, http.StatusBadRequest, err.Error())
		return
	}
	httpx.WriteJSON(w, http.StatusOK, map[string]string{"status": "logged out"})
}

// @Summary Get current user
// @Description Get the current authenticated user
// @Tags auth
// @Produce  json
// @Success 200 {object} user.User
// @Router /auth/me [get]
func (h *Handler) GetMe(w http.ResponseWriter, r *http.Request) {
	if account, ok := h.service.CurrentUser(r); ok {
		httpx.WriteJSON(w, http.StatusOK, user.SanitizeUser(*account))
		return
	}

	account, err := h.service.GetMe(extractBearerToken(r.Header.Get("Authorization")))
	if err != nil {
		httpx.WriteError(w, http.StatusUnauthorized, err.Error())
		return
	}
	httpx.WriteJSON(w, http.StatusOK, user.SanitizeUser(*account))
}

func (h *Handler) RequireAuth(next http.Handler) http.Handler {
	return h.service.Middleware(next)
}

type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}
