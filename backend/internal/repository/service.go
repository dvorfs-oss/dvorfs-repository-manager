package repository

import (
	"errors"
	"strings"

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
	db *gorm.DB
}

func NewService(db *gorm.DB) Service {
	return &service{db: db}
}

var (
	ErrRepositoryNotFound = errors.New("repository not found")
	ErrRepositoryExists   = errors.New("repository already exists")
	ErrArtifactPathEmpty  = errors.New("artifact path is required")
)

func (s *service) GetAllRepositories() ([]Repository, error) {
	var repos []Repository
	if err := s.db.Order("created_at desc").Find(&repos).Error; err != nil {
		return nil, err
	}
	return repos, nil
}

func (s *service) CreateRepository(repo *Repository) error {
	if repo.ID == uuid.Nil {
		repo.ID = uuid.New()
	}
	err := s.db.Create(repo).Error
	if err != nil && isUniqueConstraintError(err) {
		return ErrRepositoryExists
	}
	return err
}

func (s *service) GetRepository(name string) (*Repository, error) {
	var repo Repository
	if err := s.db.Where("name = ?", name).First(&repo).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrRepositoryNotFound
		}
		return nil, err
	}
	return &repo, nil
}

func (s *service) UpdateRepository(repo *Repository) error {
	var existing Repository
	if err := s.db.Where("name = ?", repo.Name).First(&existing).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return ErrRepositoryNotFound
		}
		return err
	}

	repo.ID = existing.ID
	err := s.db.Model(&existing).Updates(map[string]any{
		"format":            repo.Format,
		"type":              repo.Type,
		"attributes":        repo.Attributes,
		"cleanup_policy_id": repo.CleanupPolicyID,
		"blob_store_id":     repo.BlobStoreID,
	}).Error
	if err != nil && isUniqueConstraintError(err) {
		return ErrRepositoryExists
	}
	return err
}

func (s *service) DeleteRepository(name string) error {
	result := s.db.Where("name = ?", name).Delete(&Repository{})
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return ErrRepositoryNotFound
	}
	return nil
}

func (s *service) HandleArtifact(repoName, path string) error {
	if strings.TrimSpace(path) == "" {
		return ErrArtifactPathEmpty
	}

	repo, err := s.GetRepository(repoName)
	if err != nil {
		return err
	}

	artifact := Artifact{
		ID:           uuid.New(),
		RepositoryID: repo.ID,
		Path:         path,
	}
	return s.db.Create(&artifact).Error
}

func (s *service) SearchArtifacts(query string) ([]Artifact, error) {
	var artifacts []Artifact
	dbQuery := s.db.Order("created_at desc")
	if strings.TrimSpace(query) != "" {
		dbQuery = dbQuery.Where("LOWER(path) LIKE ?", "%"+strings.ToLower(query)+"%")
	}
	if err := dbQuery.Find(&artifacts).Error; err != nil {
		return nil, err
	}
	return artifacts, nil
}

func isUniqueConstraintError(err error) bool {
	errorText := strings.ToLower(err.Error())
	return strings.Contains(errorText, "unique") || strings.Contains(errorText, "duplicate")
}
