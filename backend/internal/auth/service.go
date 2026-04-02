package auth

import (
	"context"
	"errors"
	"net/http"
	"os"
	"strings"
	"sync"

	"dvorfs-repository-manager/internal/user"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/datatypes"
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
	db       *gorm.DB
	mu       sync.RWMutex
	sessions map[string]uuid.UUID
}

func NewService(db *gorm.DB) Service {
	svc := &service{
		db:       db,
		sessions: make(map[string]uuid.UUID),
	}
	svc.ensureBootstrapData()
	return svc
}

var (
	ErrInvalidCredentials = errors.New("invalid credentials")
	ErrInvalidToken       = errors.New("invalid token")
)

func (s *service) Login(username, password string) (string, error) {
	var account user.User
	if err := s.db.Preload("Roles").First(&account, "username = ?", strings.TrimSpace(username)).Error; err != nil {
		return "", err
	}

	if err := bcrypt.CompareHashAndPassword([]byte(account.PasswordHash), []byte(password)); err != nil {
		return "", errors.New("invalid username or password")
	}

	token := uuid.NewString()

	s.mu.Lock()
	s.sessions[token] = account.ID
	s.mu.Unlock()

	return token, nil
}

func (s *service) Logout(token string) error {
	token = strings.TrimSpace(token)
	if token == "" {
		return errors.New("token is required")
	}

	s.mu.Lock()
	delete(s.sessions, token)
	s.mu.Unlock()
	return nil
}

func (s *service) GetMe(token string) (*user.User, error) {
	token = strings.TrimSpace(token)
	if token == "" {
		return nil, errors.New("token is required")
	}

	s.mu.RLock()
	userID, ok := s.sessions[token]
	s.mu.RUnlock()
	if !ok {
		return nil, errors.New("invalid session")
	}

	var account user.User
	if err := s.db.Preload("Roles").First(&account, "id = ?", userID).Error; err != nil {
		return nil, err
	}

	return &account, nil
}

func (s *service) Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		account, err := s.GetMe(extractBearerToken(r.Header.Get("Authorization")))
		if err != nil {
			http.Error(w, "unauthorized", http.StatusUnauthorized)
			return
		}

		ctx := context.WithValue(r.Context(), userContextKey, account)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func (s *service) CurrentUser(r *http.Request) (*user.User, bool) {
	account, ok := r.Context().Value(userContextKey).(*user.User)
	return account, ok
}

func (s *service) ensureBootstrapData() {
	var adminRole user.Role
	err := s.db.First(&adminRole, "name = ?", "admin").Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		adminRole = user.Role{
			ID:         uuid.New(),
			Name:       "admin",
			Privileges: datatypes.JSON([]byte(`["*"]`)),
		}
		_ = s.db.Create(&adminRole).Error
	}

	var userCount int64
	if err := s.db.Model(&user.User{}).Count(&userCount).Error; err != nil || userCount > 0 {
		return
	}

	username := strings.TrimSpace(os.Getenv("DEFAULT_ADMIN_USERNAME"))
	if username == "" {
		username = "admin"
	}

	password := os.Getenv("DEFAULT_ADMIN_PASSWORD")
	if strings.TrimSpace(password) == "" {
		password = "admin123"
	}

	email := strings.TrimSpace(os.Getenv("DEFAULT_ADMIN_EMAIL"))
	if email == "" {
		email = "admin@local"
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return
	}

	adminUser := user.User{
		ID:           uuid.New(),
		Username:     username,
		PasswordHash: string(hashedPassword),
		Email:        email,
		Roles:        []user.Role{adminRole},
	}
	_ = s.db.Create(&adminUser).Error
}

func extractBearerToken(header string) string {
	parts := strings.SplitN(strings.TrimSpace(header), " ", 2)
	if len(parts) == 2 && strings.EqualFold(parts[0], "bearer") {
		return strings.TrimSpace(parts[1])
	}
	return strings.TrimSpace(header)
}
