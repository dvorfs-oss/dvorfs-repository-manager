package auth

import (
	"dvorfs-repository-manager/internal/user"
	"errors"
	"strings"

	"gorm.io/gorm"
)

type Service interface {
	Login(username, password string) (string, error)
	Logout(token string) error
	GetMe(token string) (*user.User, error)
}

type service struct {
	db *gorm.DB
}

func NewService(db *gorm.DB) Service {
	return &service{db: db}
}

var (
	ErrInvalidCredentials = errors.New("invalid credentials")
	ErrInvalidToken       = errors.New("invalid token")
)

func (s *service) Login(username, password string) (string, error) {
	var existingUser user.User
	if err := s.db.Where("username = ?", username).First(&existingUser).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return "", ErrInvalidCredentials
		}
		return "", err
	}
	if existingUser.PasswordHash != password {
		return "", ErrInvalidCredentials
	}
	return "token:" + username, nil
}

func (s *service) Logout(token string) error {
	if !strings.HasPrefix(token, "token:") {
		return ErrInvalidToken
	}
	return nil
}

func (s *service) GetMe(token string) (*user.User, error) {
	if !strings.HasPrefix(token, "token:") {
		return nil, ErrInvalidToken
	}
	username := strings.TrimPrefix(token, "token:")
	var existingUser user.User
	if err := s.db.Where("username = ?", username).First(&existingUser).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrInvalidToken
		}
		return nil, err
	}
	return &existingUser, nil
}
