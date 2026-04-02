package auth

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"

	"dvorfs-repository-manager/internal/user"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockService struct {
	mock.Mock
}

func (m *MockService) Login(username, password string) (string, error) {
	args := m.Called(username, password)
	return args.String(0), args.Error(1)
}

func (m *MockService) Logout(token string) error {
	args := m.Called(token)
	return args.Error(0)
}

func (m *MockService) GetMe(token string) (*user.User, error) {
	args := m.Called(token)
	return args.Get(0).(*user.User), args.Error(1)
}

func (m *MockService) Middleware(next http.Handler) http.Handler {
	return next
}

func (m *MockService) CurrentUser(r *http.Request) (*user.User, bool) {
	return nil, false
}

func TestLogin(t *testing.T) {
	mockService := new(MockService)
	handler := NewHandler(mockService)

	mockService.On("Login", "test", "test").Return("test-token", nil)
	mockService.On("GetMe", "test-token").Return(&user.User{Username: "test"}, nil)

	req, err := http.NewRequest("POST", "/api/v1/auth/login", bytes.NewBufferString(`{"username":"test","password":"test"}`))
	assert.NoError(t, err)
	req.Header.Set("Content-Type", "application/json")

	rr := httptest.NewRecorder()
	handler.Login(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)
	assert.Contains(t, rr.Body.String(), `"token":"test-token"`)
	assert.Contains(t, rr.Body.String(), `"username":"test"`)
	mockService.AssertExpectations(t)
}
