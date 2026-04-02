package user

import (
	"net/http"

	"dvorfs-repository-manager/pkg/httpx"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"gorm.io/datatypes"
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
		httpx.WriteError(w, http.StatusInternalServerError, err.Error())
		return
	}
	httpx.WriteJSON(w, http.StatusOK, SanitizeUsers(users))
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
	var request createUserRequest
	if err := httpx.DecodeJSON(r, &request); err != nil {
		httpx.WriteError(w, http.StatusBadRequest, err.Error())
		return
	}

	newUser := &User{
		Username:     request.Username,
		Email:        request.Email,
		PasswordHash: request.Password,
		Roles:        request.toRoles(),
	}
	if err := h.service.CreateUser(newUser); err != nil {
		httpx.WriteError(w, http.StatusBadRequest, err.Error())
		return
	}

	httpx.WriteJSON(w, http.StatusCreated, SanitizeUser(*newUser))
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
	var request updateUserRequest
	if err := httpx.DecodeJSON(r, &request); err != nil {
		httpx.WriteError(w, http.StatusBadRequest, err.Error())
		return
	}

	model := &User{
		Username: mux.Vars(r)["username"],
		Email:    request.Email,
	}
	if request.RoleIDs != nil {
		model.Roles = request.toRoles()
	}

	if err := h.service.UpdateUser(model); err != nil {
		httpx.WriteError(w, http.StatusBadRequest, err.Error())
		return
	}

	updated, err := h.service.GetByUsername(model.Username)
	if err != nil {
		httpx.WriteError(w, http.StatusInternalServerError, err.Error())
		return
	}

	httpx.WriteJSON(w, http.StatusOK, SanitizeUser(*updated))
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
	var request changePasswordRequest
	if err := httpx.DecodeJSON(r, &request); err != nil {
		httpx.WriteError(w, http.StatusBadRequest, err.Error())
		return
	}

	if err := h.service.ChangeUserPassword(mux.Vars(r)["username"], request.Password); err != nil {
		httpx.WriteError(w, http.StatusBadRequest, err.Error())
		return
	}

	httpx.WriteJSON(w, http.StatusOK, map[string]string{"status": "password updated"})
}

// @Summary Delete a user
// @Description Delete a user by username
// @Tags security
// @Param   username path string true "Username"
// @Success 200
// @Router /security/users/{username} [delete]
func (h *Handler) DeleteUser(w http.ResponseWriter, r *http.Request) {
	if err := h.service.DeleteUser(mux.Vars(r)["username"]); err != nil {
		httpx.WriteError(w, http.StatusInternalServerError, err.Error())
		return
	}
	httpx.WriteJSON(w, http.StatusOK, map[string]string{"status": "deleted"})
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
		httpx.WriteError(w, http.StatusInternalServerError, err.Error())
		return
	}
	httpx.WriteJSON(w, http.StatusOK, SanitizeRoles(roles))
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
	var request roleRequest
	if err := httpx.DecodeJSON(r, &request); err != nil {
		httpx.WriteError(w, http.StatusBadRequest, err.Error())
		return
	}

	role := &Role{
		Name:       request.Name,
		Privileges: request.Privileges,
	}
	if err := h.service.CreateRole(role); err != nil {
		httpx.WriteError(w, http.StatusBadRequest, err.Error())
		return
	}

	httpx.WriteJSON(w, http.StatusCreated, SanitizeRole(*role))
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
	var request roleRequest
	if err := httpx.DecodeJSON(r, &request); err != nil {
		httpx.WriteError(w, http.StatusBadRequest, err.Error())
		return
	}

	roleID, err := uuid.Parse(mux.Vars(r)["roleId"])
	if err != nil {
		httpx.WriteError(w, http.StatusBadRequest, "invalid role id")
		return
	}

	role := &Role{
		ID:         roleID,
		Name:       request.Name,
		Privileges: request.Privileges,
	}
	if err := h.service.UpdateRole(role); err != nil {
		httpx.WriteError(w, http.StatusBadRequest, err.Error())
		return
	}

	httpx.WriteJSON(w, http.StatusOK, SanitizeRole(*role))
}

// @Summary Delete a role
// @Description Delete a role by ID
// @Tags security
// @Param   roleId path string true "Role ID"
// @Success 200
// @Router /security/roles/{roleId} [delete]
func (h *Handler) DeleteRole(w http.ResponseWriter, r *http.Request) {
	if err := h.service.DeleteRole(mux.Vars(r)["roleId"]); err != nil {
		httpx.WriteError(w, http.StatusBadRequest, err.Error())
		return
	}
	httpx.WriteJSON(w, http.StatusOK, map[string]string{"status": "deleted"})
}

type createUserRequest struct {
	Username string   `json:"username"`
	Email    string   `json:"email"`
	Password string   `json:"password"`
	RoleIDs  []string `json:"roleIds"`
}

type updateUserRequest struct {
	Email   string   `json:"email"`
	RoleIDs []string `json:"roleIds"`
}

type changePasswordRequest struct {
	Password string `json:"password"`
}

type roleRequest struct {
	Name       string         `json:"name"`
	Privileges datatypes.JSON `json:"privileges"`
}

func (r createUserRequest) toRoles() []Role {
	return roleIDsToModels(r.RoleIDs)
}

func (r updateUserRequest) toRoles() []Role {
	return roleIDsToModels(r.RoleIDs)
}

func roleIDsToModels(roleIDs []string) []Role {
	roles := make([]Role, 0, len(roleIDs))
	for _, rawID := range roleIDs {
		id, err := uuid.Parse(rawID)
		if err != nil {
			continue
		}
		roles = append(roles, Role{ID: id})
	}
	return roles
}
