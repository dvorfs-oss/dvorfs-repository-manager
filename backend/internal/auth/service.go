package auth

import (
	"crypto/hmac"
	"crypto/sha256"
	"dvorfs-repository-manager/internal/user"
	"encoding/base64"
	"errors"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type contextKey string

const userContextKey contextKey = "current-user"

type Service interface {
	Login(username, password string) (string, error)
	Logout(token string) error
	GetMe(token string) (*user.User, error)
	Middleware(next http.Handler) http.Handler
	CurrentUser(r *http.Request) (*user.User, bool)
}

type service struct {
	db     *gorm.DB
	users  map[string]string
	secret []byte
}

func NewService(db *gorm.DB) Service {
	secret := os.Getenv("JWT_SECRET")
	if secret == "" {
		secret = "dev-secret"
	}

	return &service{
		db: db,
		users: map[string]string{
			"admin": "admin",
		},
		secret: []byte(secret),
	}
}

var (
	ErrInvalidCredentials = errors.New("invalid credentials")
	ErrInvalidToken       = errors.New("invalid token")
)

func (s *service) Login(username, password string) (string, error) {
	if s.db != nil {
		var foundUser user.User
		if err := s.db.Where("username = ?", username).First(&foundUser).Error; err == nil {
			if bcrypt.CompareHashAndPassword([]byte(foundUser.PasswordHash), []byte(password)) != nil {
				return "", ErrInvalidCredentials
			}
			return s.signToken(username)
		}
	}

	if expectedPassword, ok := s.users[username]; ok && expectedPassword == password {
		return s.signToken(username)
	}

	return "", ErrInvalidCredentials
}

func (s *service) Logout(token string) error {
	if _, err := s.parseToken(token); err != nil {
		return ErrInvalidToken
	}
	return nil
}

func (s *service) GetMe(token string) (*user.User, error) {
	username, err := s.parseToken(token)
	if err != nil {
		return nil, ErrInvalidToken
	}

	if s.db != nil {
		var foundUser user.User
		if err := s.db.Where("username = ?", username).First(&foundUser).Error; err == nil {
			return &foundUser, nil
		}
	}

	return &user.User{Username: username}, nil
}

func (s *service) signToken(username string) (string, error) {
	exp := time.Now().Add(24 * time.Hour).Unix()
	payload := fmt.Sprintf("%s|%d", username, exp)
	mac := hmac.New(sha256.New, s.secret)
	mac.Write([]byte(payload))
	signature := base64.RawURLEncoding.EncodeToString(mac.Sum(nil))
	token := fmt.Sprintf("%s|%s", payload, signature)
	return base64.RawURLEncoding.EncodeToString([]byte(token)), nil
}

func (s *service) parseToken(rawToken string) (string, error) {
	if strings.TrimSpace(rawToken) == "" {
		return "", ErrInvalidToken
	}

	decodedToken, err := base64.RawURLEncoding.DecodeString(rawToken)
	if err != nil {
		return "", ErrInvalidToken
	}
	parts := strings.Split(string(decodedToken), "|")
	if len(parts) != 3 {
		return "", ErrInvalidToken
	}

	username := parts[0]
	exp, err := strconv.ParseInt(parts[1], 10, 64)
	if err != nil || time.Now().Unix() > exp {
		return "", ErrInvalidToken
	}

	payload := fmt.Sprintf("%s|%d", username, exp)
	mac := hmac.New(sha256.New, s.secret)
	mac.Write([]byte(payload))
	expectedSignature := base64.RawURLEncoding.EncodeToString(mac.Sum(nil))
	if !hmac.Equal([]byte(expectedSignature), []byte(parts[2])) {
		return "", ErrInvalidToken
	}
	if strings.TrimSpace(username) == "" {
		return "", ErrInvalidToken
	}
	return username, nil
}
