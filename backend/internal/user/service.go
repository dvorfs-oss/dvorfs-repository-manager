package user

import (
	"errors"
	"strings"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type Service interface {
	GetAllUsers() ([]User, error)
	CreateUser(user *User) error
	UpdateUser(user *User) error
	ChangeUserPassword(username, newPassword string) error
	DeleteUser(username string) error
	GetAllRoles() ([]Role, error)
	CreateRole(role *Role) error
	UpdateRole(role *Role) error
	DeleteRole(roleID string) error
	GetByUsername(username string) (*User, error)
}

type service struct {
	db *gorm.DB
}

func NewService(db *gorm.DB) Service {
	return &service{db: db}
}

var (
	ErrUserNotFound = errors.New("user not found")
	ErrRoleNotFound = errors.New("role not found")
)

func (s *service) GetAllUsers() ([]User, error) {
	var users []User
	err := s.db.Preload("Roles").Order("username asc").Find(&users).Error
	return users, err
}

func (s *service) CreateUser(model *User) error {
	if model == nil {
		return errors.New("user is required")
	}

	username := strings.TrimSpace(model.Username)
	email := strings.TrimSpace(model.Email)
	password := model.PasswordHash
	if username == "" || email == "" || password == "" {
		return errors.New("username, email and password are required")
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	roles, err := s.resolveRoles(model.Roles)
	if err != nil {
		return err
	}

	model.ID = uuid.New()
	model.Username = username
	model.Email = email
	model.PasswordHash = string(hashedPassword)
	model.Roles = roles

	return s.db.Create(model).Error
}

func (s *service) UpdateUser(model *User) error {
	if model == nil {
		return errors.New("user is required")
	}

	username := strings.TrimSpace(model.Username)
	if username == "" {
		return errors.New("username is required")
	}

	var existing User
	if err := s.db.Preload("Roles").First(&existing, "username = ?", username).Error; err != nil {
		return err
	}

	email := strings.TrimSpace(model.Email)
	if email == "" {
		return errors.New("email is required")
	}

	existing.Email = email
	if model.Roles != nil {
		roles, err := s.resolveRoles(model.Roles)
		if err != nil {
			return err
		}
		if err := s.db.Model(&existing).Association("Roles").Replace(roles); err != nil {
			return err
		}
		existing.Roles = roles
	}

	return s.db.Save(&existing).Error
}

func (s *service) ChangeUserPassword(username, newPassword string) error {
	if strings.TrimSpace(newPassword) == "" {
		return errors.New("password is required")
	}

	var existing User
	if err := s.db.First(&existing, "username = ?", strings.TrimSpace(username)).Error; err != nil {
		return err
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(newPassword), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	return s.db.Model(&existing).Update("password_hash", string(hashedPassword)).Error
}

func (s *service) DeleteUser(username string) error {
	return s.db.Delete(&User{}, "username = ?", strings.TrimSpace(username)).Error
}

func (s *service) GetAllRoles() ([]Role, error) {
	var roles []Role
	err := s.db.Order("name asc").Find(&roles).Error
	return roles, err
}

func (s *service) CreateRole(role *Role) error {
	if role == nil {
		return errors.New("role is required")
	}

	role.Name = strings.TrimSpace(role.Name)
	if role.Name == "" {
		return errors.New("role name is required")
	}

	if role.ID == uuid.Nil {
		role.ID = uuid.New()
	}

	return s.db.Create(role).Error
}

func (s *service) UpdateRole(role *Role) error {
	if role == nil || role.ID == uuid.Nil {
		return errors.New("role id is required")
	}

	var existing Role
	if err := s.db.First(&existing, "id = ?", role.ID).Error; err != nil {
		return err
	}

	existing.Name = strings.TrimSpace(role.Name)
	existing.Privileges = role.Privileges
	return s.db.Save(&existing).Error
}

func (s *service) DeleteRole(roleID string) error {
	id, err := uuid.Parse(strings.TrimSpace(roleID))
	if err != nil {
		return err
	}
	return s.db.Delete(&Role{}, "id = ?", id).Error
}

func (s *service) GetByUsername(username string) (*User, error) {
	var existing User
	if err := s.db.Preload("Roles").First(&existing, "username = ?", strings.TrimSpace(username)).Error; err != nil {
		return nil, err
	}
	return &existing, nil
}

func (s *service) resolveRoles(roles []Role) ([]Role, error) {
	if len(roles) == 0 {
		return nil, nil
	}

	ids := make([]uuid.UUID, 0, len(roles))
	for _, role := range roles {
		if role.ID != uuid.Nil {
			ids = append(ids, role.ID)
		}
	}

	if len(ids) == 0 {
		return nil, nil
	}

	var resolved []Role
	if err := s.db.Where("id IN ?", ids).Find(&resolved).Error; err != nil {
		return nil, err
	}

	if len(resolved) != len(ids) {
		return nil, errors.New("one or more roles were not found")
	}

	return resolved, nil
}
