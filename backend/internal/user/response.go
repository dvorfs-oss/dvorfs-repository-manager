package user

import (
	"encoding/json"
	"time"

	"github.com/google/uuid"
)

type UserResponse struct {
	ID        uuid.UUID      `json:"id"`
	Username  string         `json:"username"`
	Email     string         `json:"email"`
	CreatedAt time.Time      `json:"createdAt"`
	UpdatedAt time.Time      `json:"updatedAt"`
	Roles     []RoleResponse `json:"roles"`
}

type RoleResponse struct {
	ID         uuid.UUID `json:"id"`
	Name       string    `json:"name"`
	Privileges any       `json:"privileges"`
	CreatedAt  time.Time `json:"createdAt"`
}

func SanitizeUser(model User) UserResponse {
	roles := make([]RoleResponse, 0, len(model.Roles))
	for _, role := range model.Roles {
		roles = append(roles, SanitizeRole(role))
	}

	return UserResponse{
		ID:        model.ID,
		Username:  model.Username,
		Email:     model.Email,
		CreatedAt: model.CreatedAt,
		UpdatedAt: model.UpdatedAt,
		Roles:     roles,
	}
}

func SanitizeUsers(models []User) []UserResponse {
	result := make([]UserResponse, 0, len(models))
	for _, item := range models {
		result = append(result, SanitizeUser(item))
	}
	return result
}

func SanitizeRole(model Role) RoleResponse {
	return RoleResponse{
		ID:         model.ID,
		Name:       model.Name,
		Privileges: decodeJSON(model.Privileges),
		CreatedAt:  model.CreatedAt,
	}
}

func SanitizeRoles(models []Role) []RoleResponse {
	result := make([]RoleResponse, 0, len(models))
	for _, item := range models {
		result = append(result, SanitizeRole(item))
	}
	return result
}

func decodeJSON(data []byte) any {
	if len(data) == 0 {
		return nil
	}

	var value any
	if err := json.Unmarshal(data, &value); err != nil {
		return string(data)
	}
	return value
}
