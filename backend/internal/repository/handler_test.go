package repository

import (
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockService struct {
	mock.Mock
}

func (m *MockService) GetAllRepositories() ([]Repository, error) {
	args := m.Called()
	return args.Get(0).([]Repository), args.Error(1)
}

func (m *MockService) CreateRepository(repo *Repository) error {
	args := m.Called(repo)
	return args.Error(0)
}

func (m *MockService) GetRepository(name string) (*Repository, error) {
	args := m.Called(name)
	return args.Get(0).(*Repository), args.Error(1)
}

func (m *MockService) UpdateRepository(repo *Repository) error {
	args := m.Called(repo)
	return args.Error(0)
}

func (m *MockService) DeleteRepository(name string) error {
	args := m.Called(name)
	return args.Error(0)
}

func (m *MockService) UploadArtifact(repoName, artifactPath, contentType string, body io.Reader) (*Artifact, error) {
	args := m.Called(repoName, artifactPath, contentType, body)
	artifact, _ := args.Get(0).(*Artifact)
	return artifact, args.Error(1)
}

func (m *MockService) OpenArtifact(repoName, artifactPath string) (io.ReadCloser, *Artifact, error) {
	args := m.Called(repoName, artifactPath)
	reader, _ := args.Get(0).(io.ReadCloser)
	artifact, _ := args.Get(1).(*Artifact)
	return reader, artifact, args.Error(2)
}

func (m *MockService) DeleteArtifact(repoName, artifactPath string) error {
	args := m.Called(repoName, artifactPath)
	return args.Error(0)
}

func (m *MockService) ListArtifacts(repoName string) ([]Artifact, error) {
	args := m.Called(repoName)
	return args.Get(0).([]Artifact), args.Error(1)
}

func (m *MockService) SearchArtifacts(query string) ([]Artifact, error) {
	args := m.Called(query)
	return args.Get(0).([]Artifact), args.Error(1)
}

func TestGetAllRepositories(t *testing.T) {
	mockService := new(MockService)
	handler := NewHandler(mockService)

	mockService.On("GetAllRepositories").Return([]Repository{{Name: "test-repo"}}, nil)

	req, err := http.NewRequest("GET", "/api/v1/repositories", nil)
	assert.NoError(t, err)

	rr := httptest.NewRecorder()
	handler.GetAllRepositories(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)
	assert.JSONEq(t, `[{"ID":"00000000-0000-0000-0000-000000000000","Name":"test-repo","Format":"","Type":"","Attributes":null,"CleanupPolicyID":null,"BlobStoreID":null,"CreatedAt":"0001-01-01T00:00:00Z","UpdatedAt":"0001-01-01T00:00:00Z","Artifacts":null,"CleanupPolicy":null,"BlobStore":null}]`, rr.Body.String())
	mockService.AssertExpectations(t)
}
