package repository

import (
	"github.com/google/uuid"
)

type Service interface {
	GetAllRepositories() ([]Repository, error)
	CreateRepository(repo *Repository) error
	GetRepository(name string) (*Repository, error)
	UpdateRepository(repo *Repository) error
	DeleteRepository(name string) error
	HandleArtifact(repoName, path string) error
	SearchArtifacts(query string) ([]Artifact, error)
}

type service struct {
	// Add dependencies here, e.g., a repository repository
}

func NewService() Service {
	return &service{}
}

func (s *service) GetAllRepositories() ([]Repository, error) {
	return []Repository{{Name: "test-repo"}}, nil
}

func (s *service) CreateRepository(repo *Repository) error {
	return nil
}

func (s *service) GetRepository(name string) (*Repository, error) {
	return &Repository{Name: name}, nil
}

func (s *service) UpdateRepository(repo *Repository) error {
	return nil
}

func (s *service) DeleteRepository(name string) error {
	return nil
}

func (s *service) HandleArtifact(repoName, path string) error {
	return nil
}

func (s *service) SearchArtifacts(query string) ([]Artifact, error) {
	return []Artifact{{ID: uuid.New()}}, nil
}
