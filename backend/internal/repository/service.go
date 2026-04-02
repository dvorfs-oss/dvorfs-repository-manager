package repository

import (
	"errors"
	"strings"
	"sync"

	"github.com/google/uuid"
	"gorm.io/gorm"
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
	mu          sync.RWMutex
	reposByName map[string]Repository
	artifacts   []Artifact
	db          *gorm.DB
}

func NewService(db *gorm.DB) Service {
	return &service{
		reposByName: make(map[string]Repository),
		artifacts:   make([]Artifact, 0),
		db:          db,
	}
}

var (
	ErrRepositoryNotFound = errors.New("repository not found")
	ErrRepositoryExists   = errors.New("repository already exists")
	ErrArtifactPathEmpty  = errors.New("artifact path is required")
)

func (s *service) GetAllRepositories() ([]Repository, error) {
	if s.db != nil {
		var repos []Repository
		if err := s.db.Find(&repos).Error; err != nil {
			return nil, err
		}
		return repos, nil
	}

	s.mu.RLock()
	defer s.mu.RUnlock()

	repos := make([]Repository, 0, len(s.reposByName))
	for _, repo := range s.reposByName {
		repos = append(repos, repo)
	}

	return repos, nil
}

func (s *service) CreateRepository(repo *Repository) error {
	if s.db != nil {
		repo.ID = uuid.New()
		if err := s.db.Create(repo).Error; err != nil {
			if strings.Contains(strings.ToLower(err.Error()), "duplicate") {
				return ErrRepositoryExists
			}
			return err
		}
		return nil
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	if _, exists := s.reposByName[repo.Name]; exists {
		return ErrRepositoryExists
	}
	if repo.ID == uuid.Nil {
		repo.ID = uuid.New()
	}
	s.reposByName[repo.Name] = *repo
	return nil
}

func (s *service) GetRepository(name string) (*Repository, error) {
	if s.db != nil {
		var repo Repository
		if err := s.db.Where("name = ?", name).First(&repo).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return nil, ErrRepositoryNotFound
			}
			return nil, err
		}
		return &repo, nil
	}

	s.mu.RLock()
	defer s.mu.RUnlock()

	repo, exists := s.reposByName[name]
	if !exists {
		return nil, ErrRepositoryNotFound
	}
	return &repo, nil
}

func (s *service) UpdateRepository(repo *Repository) error {
	if s.db != nil {
		var existing Repository
		if err := s.db.Where("name = ?", repo.Name).First(&existing).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return ErrRepositoryNotFound
			}
			return err
		}
		repo.ID = existing.ID
		return s.db.Save(repo).Error
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	if _, exists := s.reposByName[repo.Name]; !exists {
		return ErrRepositoryNotFound
	}
	s.reposByName[repo.Name] = *repo
	return nil
}

func (s *service) DeleteRepository(name string) error {
	if s.db != nil {
		result := s.db.Where("name = ?", name).Delete(&Repository{})
		if result.Error != nil {
			return result.Error
		}
		if result.RowsAffected == 0 {
			return ErrRepositoryNotFound
		}
		return nil
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	if _, exists := s.reposByName[name]; !exists {
		return ErrRepositoryNotFound
	}
	delete(s.reposByName, name)
	return nil
}

func (s *service) HandleArtifact(repoName, path string) error {
	if s.db != nil {
		var repo Repository
		if err := s.db.Where("name = ?", repoName).First(&repo).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return ErrRepositoryNotFound
			}
			return err
		}
		if strings.TrimSpace(path) == "" {
			return ErrArtifactPathEmpty
		}
		artifact := Artifact{
			ID:           uuid.New(),
			RepositoryID: repo.ID,
			Path:         path,
		}
		return s.db.Create(&artifact).Error
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	if _, exists := s.reposByName[repoName]; !exists {
		return ErrRepositoryNotFound
	}
	if strings.TrimSpace(path) == "" {
		return ErrArtifactPathEmpty
	}

	s.artifacts = append(s.artifacts, Artifact{
		ID:           uuid.New(),
		RepositoryID: s.reposByName[repoName].ID,
		Path:         path,
	})
	return nil
}

func (s *service) SearchArtifacts(query string) ([]Artifact, error) {
	if s.db != nil {
		var artifacts []Artifact
		tx := s.db.Model(&Artifact{})
		if strings.TrimSpace(query) != "" {
			tx = tx.Where("LOWER(path) LIKE ?", "%"+strings.ToLower(query)+"%")
		}
		if err := tx.Find(&artifacts).Error; err != nil {
			return nil, err
		}
		return artifacts, nil
	}

	s.mu.RLock()
	defer s.mu.RUnlock()

	if strings.TrimSpace(query) == "" {
		return s.artifacts, nil
	}

	query = strings.ToLower(query)
	result := make([]Artifact, 0)
	for _, artifact := range s.artifacts {
		if strings.Contains(strings.ToLower(artifact.Path), query) {
			result = append(result, artifact)
		}
	}
	return result, nil
}
