package auth

import (
	"dvorfs-repository-manager/internal/user"
	"net/http"
	"net/http/httptest"
	"testing"

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

func TestLogin(t *testing.T) {
	mockService := new(MockService)
	handler := NewHandler(mockService)

	mockService.On("Login", "test", "test").Return("test-token", nil)

	req, err := http.NewRequest("POST", "/api/v1/auth/login", nil)
	assert.NoError(t, err)

	rr := httptest.NewRecorder()
	handler.Login(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)
	assert.JSONEq(t, `{"token":"test-token"}`, rr.Body.String())
	mockService.AssertExpectations(t)
}
