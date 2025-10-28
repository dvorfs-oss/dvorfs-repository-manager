package auth

import (
	"dvorfs-repository-manager/internal/user"
)

type Service interface {
	Login(username, password string) (string, error)
	Logout(token string) error
	GetMe(token string) (*user.User, error)
}

type service struct {
	// Add dependencies here, e.g., a user repository
}

func NewService() Service {
	return &service{}
}

func (s *service) Login(username, password string) (string, error) {
	// Business logic for logging in a user
	return "jwt-token", nil
}

func (s *service) Logout(token string) error {
	// Business logic for logging out a user
	return nil
}

func (s *service) GetMe(token string) (*user.User, error) {
	// Business logic for getting the current user
	return &user.User{Username: "testuser"}, nil
}
