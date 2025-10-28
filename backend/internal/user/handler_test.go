package user

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockService struct {
	mock.Mock
}

func (m *MockService) GetAllUsers() ([]User, error) {
	args := m.Called()
	return args.Get(0).([]User), args.Error(1)
}

func (m *MockService) CreateUser(user *User) error {
	args := m.Called(user)
	return args.Error(0)
}

func (m *MockService) UpdateUser(user *User) error {
	args := m.Called(user)
	return args.Error(0)
}

func (m *MockService) ChangeUserPassword(username, newPassword string) error {
	args := m.Called(username, newPassword)
	return args.Error(0)
}

func (m *MockService) DeleteUser(username string) error {
	args := m.Called(username)
	return args.Error(0)
}

func (m *MockService) GetAllRoles() ([]Role, error) {
	args := m.Called()
	return args.Get(0).([]Role), args.Error(1)
}

func (m *MockService) CreateRole(role *Role) error {
	args := m.Called(role)
	return args.Error(0)
}

func (m *MockService) UpdateRole(role *Role) error {
	args := m.Called(role)
	return args.Error(0)
}

func (m *MockService) DeleteRole(roleID string) error {
	args := m.Called(roleID)
	return args.Error(0)
}

func TestGetAllUsers(t *testing.T) {
	mockService := new(MockService)
	handler := NewHandler(mockService)

	mockService.On("GetAllUsers").Return([]User{{Username: "testuser"}}, nil)

	req, err := http.NewRequest("GET", "/api/v1/security/users", nil)
	assert.NoError(t, err)

	rr := httptest.NewRecorder()
	handler.GetAllUsers(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)
	assert.JSONEq(t, `[{"ID":"00000000-0000-0000-0000-000000000000","Username":"testuser","PasswordHash":"","Email":"","CreatedAt":"0001-01-01T00:00:00Z","UpdatedAt":"0001-01-01T00:00:00Z","Roles":null}]`, rr.Body.String())
	mockService.AssertExpectations(t)
}
