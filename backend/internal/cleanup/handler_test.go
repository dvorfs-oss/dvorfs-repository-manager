package cleanup

import (
	"dvorfs-repository-manager/internal/repository"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockService struct {
	mock.Mock
}

func (m *MockService) GetAllCleanupPolicies() ([]repository.CleanupPolicy, error) {
	args := m.Called()
	return args.Get(0).([]repository.CleanupPolicy), args.Error(1)
}

func (m *MockService) CreateCleanupPolicy(policy *repository.CleanupPolicy) error {
	args := m.Called(policy)
	return args.Error(0)
}

func (m *MockService) UpdateCleanupPolicy(policy *repository.CleanupPolicy) error {
	args := m.Called(policy)
	return args.Error(0)
}

func (m *MockService) DeleteCleanupPolicy(policyID string) error {
	args := m.Called(policyID)
	return args.Error(0)
}

func TestGetAllCleanupPolicies(t *testing.T) {
	mockService := new(MockService)
	handler := NewHandler(mockService)

	mockService.On("GetAllCleanupPolicies").Return([]repository.CleanupPolicy{{Name: "test-policy"}}, nil)

	req, err := http.NewRequest("GET", "/api/v1/cleanup-policies", nil)
	assert.NoError(t, err)

	rr := httptest.NewRecorder()
	handler.GetAllCleanupPolicies(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)
	assert.JSONEq(t, `[{"ID":"00000000-0000-0000-0000-000000000000","Name":"test-policy","Criteria":null,"CreatedAt":"0001-01-01T00:00:00Z","UpdatedAt":"0001-01-01T00:00:00Z"}]`, rr.Body.String())
	mockService.AssertExpectations(t)
}
