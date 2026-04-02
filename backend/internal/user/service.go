package user

import (
	"errors"

	"github.com/google/uuid"
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
	if err := s.db.Preload("Roles").Order("created_at desc").Find(&users).Error; err != nil {
		return nil, err
	}
	return users, nil
}

func (s *service) CreateUser(user *User) error {
	if user.ID == uuid.Nil {
		user.ID = uuid.New()
	}
	return s.db.Create(user).Error
}

func (s *service) UpdateUser(user *User) error {
	var existing User
	if err := s.db.Where("username = ?", user.Username).First(&existing).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return ErrUserNotFound
		}
		return err
	}
	user.ID = existing.ID
	if err := s.db.Model(&existing).Updates(map[string]any{
		"email":         user.Email,
		"password_hash": user.PasswordHash,
	}).Error; err != nil {
		return err
	}
	return nil
}

func (s *service) ChangeUserPassword(username, newPassword string) error {
	result := s.db.Model(&User{}).Where("username = ?", username).Update("password_hash", newPassword)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return ErrUserNotFound
	}
	return nil
}

func (s *service) DeleteUser(username string) error {
	result := s.db.Where("username = ?", username).Delete(&User{})
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return ErrUserNotFound
	}
	return nil
}

func (s *service) GetAllRoles() ([]Role, error) {
	var roles []Role
	if err := s.db.Order("created_at desc").Find(&roles).Error; err != nil {
		return nil, err
	}
	return roles, nil
}

func (s *service) CreateRole(role *Role) error {
	if role.ID == uuid.Nil {
		role.ID = uuid.New()
	}
	return s.db.Create(role).Error
}

func (s *service) UpdateRole(role *Role) error {
	result := s.db.Model(&Role{}).Where("id = ?", role.ID).Updates(map[string]any{
		"name":       role.Name,
		"privileges": role.Privileges,
	})
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return ErrRoleNotFound
	}
	return nil
}

func (s *service) DeleteRole(roleID string) error {
	result := s.db.Where("id = ?", roleID).Delete(&Role{})
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return ErrRoleNotFound
	}
	return nil
}
