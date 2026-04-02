package user

import (
	"errors"
	"sync"

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
}

type service struct {
	mu    sync.RWMutex
	users map[string]User
	roles map[string]Role
	db    *gorm.DB
}

func NewService(db *gorm.DB) Service {
	return &service{
		users: make(map[string]User),
		roles: make(map[string]Role),
		db:    db,
	}
}

var (
	ErrUserNotFound = errors.New("user not found")
	ErrRoleNotFound = errors.New("role not found")
)

func (s *service) GetAllUsers() ([]User, error) {
	if s.db != nil {
		var users []User
		if err := s.db.Preload("Roles").Find(&users).Error; err != nil {
			return nil, err
		}
		return users, nil
	}

	s.mu.RLock()
	defer s.mu.RUnlock()
	result := make([]User, 0, len(s.users))
	for _, u := range s.users {
		result = append(result, u)
	}
	return result, nil
}

func (s *service) CreateUser(user *User) error {
	passwordHash, err := hashPassword(user.PasswordHash)
	if err != nil {
		return err
	}
	user.PasswordHash = passwordHash

	if s.db != nil {
		user.ID = uuid.New()
		return s.db.Create(user).Error
	}

	s.mu.Lock()
	defer s.mu.Unlock()
	if user.ID == uuid.Nil {
		user.ID = uuid.New()
	}
	s.users[user.Username] = *user
	return nil
}

func (s *service) UpdateUser(user *User) error {
	if s.db != nil {
		var existing User
		if err := s.db.Where("username = ?", user.Username).First(&existing).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return ErrUserNotFound
			}
			return err
		}
		user.ID = existing.ID
		return s.db.Save(user).Error
	}

	s.mu.Lock()
	defer s.mu.Unlock()
	if _, exists := s.users[user.Username]; !exists {
		return ErrUserNotFound
	}
	s.users[user.Username] = *user
	return nil
}

func (s *service) ChangeUserPassword(username, newPassword string) error {
	passwordHash, err := hashPassword(newPassword)
	if err != nil {
		return err
	}

	if s.db != nil {
		result := s.db.Model(&User{}).Where("username = ?", username).Update("password_hash", passwordHash)
		if result.Error != nil {
			return result.Error
		}
		if result.RowsAffected == 0 {
			return ErrUserNotFound
		}
		return nil
	}

	s.mu.Lock()
	defer s.mu.Unlock()
	currentUser, exists := s.users[username]
	if !exists {
		return ErrUserNotFound
	}
	currentUser.PasswordHash = passwordHash
	s.users[username] = currentUser
	return nil
}

func (s *service) DeleteUser(username string) error {
	if s.db != nil {
		result := s.db.Where("username = ?", username).Delete(&User{})
		if result.Error != nil {
			return result.Error
		}
		if result.RowsAffected == 0 {
			return ErrUserNotFound
		}
		return nil
	}

	s.mu.Lock()
	defer s.mu.Unlock()
	if _, exists := s.users[username]; !exists {
		return ErrUserNotFound
	}
	delete(s.users, username)
	return nil
}

func (s *service) GetAllRoles() ([]Role, error) {
	if s.db != nil {
		var roles []Role
		if err := s.db.Find(&roles).Error; err != nil {
			return nil, err
		}
		return roles, nil
	}

	s.mu.RLock()
	defer s.mu.RUnlock()
	result := make([]Role, 0, len(s.roles))
	for _, role := range s.roles {
		result = append(result, role)
	}
	return result, nil
}

func (s *service) CreateRole(role *Role) error {
	if s.db != nil {
		role.ID = uuid.New()
		return s.db.Create(role).Error
	}

	s.mu.Lock()
	defer s.mu.Unlock()
	if role.ID == uuid.Nil {
		role.ID = uuid.New()
	}
	s.roles[role.ID.String()] = *role
	return nil
}

func (s *service) UpdateRole(role *Role) error {
	if s.db != nil {
		var existing Role
		if err := s.db.First(&existing, "id = ?", role.ID).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return ErrRoleNotFound
			}
			return err
		}
		return s.db.Save(role).Error
	}

	s.mu.Lock()
	defer s.mu.Unlock()
	if _, exists := s.roles[role.ID.String()]; !exists {
		return ErrRoleNotFound
	}
	s.roles[role.ID.String()] = *role
	return nil
}

func (s *service) DeleteRole(roleID string) error {
	if s.db != nil {
		result := s.db.Delete(&Role{}, "id = ?", roleID)
		if result.Error != nil {
			return result.Error
		}
		if result.RowsAffected == 0 {
			return ErrRoleNotFound
		}
		return nil
	}

	s.mu.Lock()
	defer s.mu.Unlock()
	if _, exists := s.roles[roleID]; !exists {
		return ErrRoleNotFound
	}
	delete(s.roles, roleID)
	return nil
}

func hashPassword(rawPassword string) (string, error) {
	if rawPassword == "" {
		return "", errors.New("password is required")
	}
	passwordHash, err := bcrypt.GenerateFromPassword([]byte(rawPassword), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(passwordHash), nil
}
