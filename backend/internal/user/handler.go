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

// @Summary Get all users
// @Description Get a list of all users
// @Tags security
// @Produce  json
// @Success 200 {array} User
// @Router /security/users [get]
func (h *Handler) GetAllUsers(w http.ResponseWriter, r *http.Request) {
	users, err := h.service.GetAllUsers()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(users)
}

// @Summary Create a new user
// @Description Create a new user
// @Tags security
// @Accept  json
// @Produce  json
// @Param   user body User true "User"
// @Success 201
// @Router /security/users [post]
func (h *Handler) CreateUser(w http.ResponseWriter, r *http.Request) {
	// Implementation needed
	w.WriteHeader(http.StatusCreated)
}

// @Summary Update a user
// @Description Update a user by username
// @Tags security
// @Accept  json
// @Produce  json
// @Param   username path string true "Username"
// @Param   user body User true "User"
// @Success 200
// @Router /security/users/{username} [put]
func (h *Handler) UpdateUser(w http.ResponseWriter, r *http.Request) {
	// Implementation needed
	w.WriteHeader(http.StatusOK)
}

// @Summary Change user password
// @Description Change a user's password
// @Tags security
// @Accept  json
// @Produce  json
// @Param   username path string true "Username"
// @Param   password body map[string]string true "Password"
// @Success 200
// @Router /security/users/{username}/password [put]
func (h *Handler) ChangeUserPassword(w http.ResponseWriter, r *http.Request) {
	// Implementation needed
	w.WriteHeader(http.StatusOK)
}

// @Summary Delete a user
// @Description Delete a user by username
// @Tags security
// @Param   username path string true "Username"
// @Success 200
// @Router /security/users/{username} [delete]
func (h *Handler) DeleteUser(w http.ResponseWriter, r *http.Request) {
	// Implementation needed
	w.WriteHeader(http.StatusOK)
}

// @Summary Get all roles
// @Description Get a list of all roles
// @Tags security
// @Produce  json
// @Success 200 {array} Role
// @Router /security/roles [get]
func (h *Handler) GetAllRoles(w http.ResponseWriter, r *http.Request) {
	roles, err := h.service.GetAllRoles()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(roles)
}

// @Summary Create a new role
// @Description Create a new role
// @Tags security
// @Accept  json
// @Produce  json
// @Param   role body Role true "Role"
// @Success 201
// @Router /security/roles [post]
func (h *Handler) CreateRole(w http.ResponseWriter, r *http.Request) {
	// Implementation needed
	w.WriteHeader(http.StatusCreated)
}

// @Summary Update a role
// @Description Update a role by ID
// @Tags security
// @Accept  json
// @Produce  json
// @Param   roleId path string true "Role ID"
// @Param   role body Role true "Role"
// @Success 200
// @Router /security/roles/{roleId} [put]
func (h *Handler) UpdateRole(w http.ResponseWriter, r *http.Request) {
	// Implementation needed
	w.WriteHeader(http.StatusOK)
}

// @Summary Delete a role
// @Description Delete a role by ID
// @Tags security
// @Param   roleId path string true "Role ID"
// @Success 200
// @Router /security/roles/{roleId} [delete]
func (h *Handler) DeleteRole(w http.ResponseWriter, r *http.Request) {
	// Implementation needed
	w.WriteHeader(http.StatusOK)
}
