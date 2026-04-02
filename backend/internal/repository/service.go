package repository

import (
	"errors"
	"fmt"
	"io"
	"log"
	"os"
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
	UploadArtifact(repoName, artifactPath, contentType string, body io.Reader) (*Artifact, error)
	OpenArtifact(repoName, artifactPath string) (io.ReadCloser, *Artifact, error)
	DeleteArtifact(repoName, artifactPath string) error
	ListArtifacts(repoName string) ([]Artifact, error)
	SearchArtifacts(query string) ([]Artifact, error)
}

type service struct {
	db     *gorm.DB
	storer *artifactStorage
}

func NewService(db *gorm.DB) Service {
	storer, err := newArtifactStorage(defaultArtifactStorageRoot())
	if err != nil {
		log.Fatal(err)
	}

	return &service{
		db:     db,
		storer: storer,
	}
}

func (s *service) GetAllRepositories() ([]Repository, error) {
	var repos []Repository
	if err := s.db.Order("created_at DESC").Find(&repos).Error; err != nil {
		return nil, err
	}
	return repos, nil
}

func (s *service) CreateRepository(repo *Repository) error {
	if repo == nil {
		return errors.New("repository is required")
	}

	normalizeRepository(repo)
	if repo.Name == "" {
		return errors.New("repository name is required")
	}

	if repo.Type == "" {
		repo.Type = "hosted"
	}
	if repo.Format == "" {
		repo.Format = "raw"
	}
	if repo.ID == uuid.Nil {
		repo.ID = uuid.New()
	}

	if err := s.storer.ensureRepositoryRoot(repo.Name); err != nil {
		return err
	}
	if err := s.db.Create(repo).Error; err != nil {
		return err
	}
	return nil
}

func (s *service) GetRepository(name string) (*Repository, error) {
	var repo Repository
	if err := s.db.Preload("Artifacts").Preload("CleanupPolicy").Preload("BlobStore").First(&repo, "name = ?", name).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("%w: %s", ErrRepositoryNotFound, name)
		}
		return nil, err
	}
	return &repo, nil
}

func (s *service) UpdateRepository(repo *Repository) error {
	if repo == nil {
		return errors.New("repository is required")
	}

	normalizeRepository(repo)
	if repo.Name == "" {
		return errors.New("repository name is required")
	}

	updates := map[string]any{
		"format":            repo.Format,
		"type":              repo.Type,
		"attributes":        repo.Attributes,
		"cleanup_policy_id": repo.CleanupPolicyID,
		"blob_store_id":     repo.BlobStoreID,
	}

	result := s.db.Model(&Repository{}).Where("name = ?", repo.Name).Updates(updates)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return fmt.Errorf("%w: %s", ErrRepositoryNotFound, repo.Name)
	}

	return nil
}

func (s *service) DeleteRepository(name string) error {
	repo, err := s.GetRepository(name)
	if err != nil {
		return err
	}

	if err := s.db.Where("repository_id = ?", repo.ID).Delete(&Artifact{}).Error; err != nil {
		return err
	}
	if err := s.db.Delete(&Repository{}, "id = ?", repo.ID).Error; err != nil {
		return err
	}
	return s.storer.deleteRepositoryRoot(name)
}

func (s *service) UploadArtifact(repoName, artifactPath, contentType string, body io.Reader) (*Artifact, error) {
	repo, err := s.GetRepository(repoName)
	if err != nil {
		return nil, err
	}
	if !strings.EqualFold(repo.Type, "hosted") {
		return nil, fmt.Errorf("%w: %s", ErrUnsupportedRepositoryType, repo.Type)
	}
	normalizedPath, err := normalizeArtifactPathStrict(artifactPath)
	if err != nil {
		return nil, err
	}

	fullPath, size, checksums, err := s.storer.saveArtifact(repoName, normalizedPath, body)
	if err != nil {
		return nil, err
	}

	artifact := &Artifact{
		RepositoryID: repo.ID,
		Path:         normalizedPath,
		Size:         size,
		ContentType:  strings.TrimSpace(contentType),
		Checksums:    checksums.toJSON(),
	}

	var existing Artifact
	result := s.db.Where("repository_id = ? AND path = ?", repo.ID, artifact.Path).First(&existing)
	if result.Error != nil && !errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil, result.Error
	}

	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		artifact.ID = uuid.New()
		if err := s.db.Create(artifact).Error; err != nil {
			_ = os.Remove(fullPath)
			return nil, err
		}
		return artifact, nil
	}

	existing.Size = size
	existing.ContentType = artifact.ContentType
	existing.Checksums = artifact.Checksums
	if err := s.db.Save(&existing).Error; err != nil {
		_ = os.Remove(fullPath)
		return nil, err
	}
	return &existing, nil
}

func (s *service) SearchArtifacts(query string) ([]Artifact, error) {
	var artifacts []Artifact
	trimmed := strings.TrimSpace(query)
	db := s.db.Order("created_at DESC")
	if trimmed == "" {
		if err := db.Find(&artifacts).Error; err != nil {
			return nil, err
		}
		return artifacts, nil
	}

	like := "%" + strings.ToLower(trimmed) + "%"
	if err := db.Where("LOWER(path) LIKE ? OR LOWER(content_type) LIKE ?", like, like).Find(&artifacts).Error; err != nil {
		return nil, err
	}
	return artifacts, nil
}

func (s *service) OpenArtifact(repoName, artifactPath string) (io.ReadCloser, *Artifact, error) {
	repo, err := s.GetRepository(repoName)
	if err != nil {
		return nil, nil, err
	}
	normalizedPath, err := normalizeArtifactPathStrict(artifactPath)
	if err != nil {
		return nil, nil, err
	}

	var artifact Artifact
	if err := s.db.Where("repository_id = ? AND path = ?", repo.ID, normalizedPath).First(&artifact).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil, fmt.Errorf("%w: %s", ErrArtifactNotFound, artifactPath)
		}
		return nil, nil, err
	}

	file, err := s.storer.openArtifact(repoName, normalizedPath)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return nil, nil, fmt.Errorf("%w: %s", ErrArtifactNotFound, artifactPath)
		}
		return nil, nil, err
	}

	now := currentTime()
	artifact.LastDownloadedAt = &now
	if err := s.db.Model(&Artifact{}).Where("id = ?", artifact.ID).Update("last_downloaded_at", now).Error; err != nil {
		_ = file.Close()
		return nil, nil, err
	}

	return file, &artifact, nil
}

func (s *service) DeleteArtifact(repoName, artifactPath string) error {
	repo, err := s.GetRepository(repoName)
	if err != nil {
		return err
	}
	normalizedPath, err := normalizeArtifactPathStrict(artifactPath)
	if err != nil {
		return err
	}

	var artifact Artifact
	if err := s.db.Where("repository_id = ? AND path = ?", repo.ID, normalizedPath).First(&artifact).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return fmt.Errorf("%w: %s", ErrArtifactNotFound, artifactPath)
		}
		return err
	}

	if err := s.storer.deleteArtifact(repoName, normalizedPath); err != nil && !errors.Is(err, os.ErrNotExist) {
		return err
	}

	return s.db.Delete(&Artifact{}, "id = ?", artifact.ID).Error
}

func (s *service) ListArtifacts(repoName string) ([]Artifact, error) {
	repo, err := s.GetRepository(repoName)
	if err != nil {
		return nil, err
	}

	var artifacts []Artifact
	if err := s.db.Where("repository_id = ?", repo.ID).Order("path ASC").Find(&artifacts).Error; err != nil {
		return nil, err
	}
	return artifacts, nil
}

func normalizeRepository(repo *Repository) {
	repo.Name = strings.TrimSpace(repo.Name)
	repo.Format = strings.ToLower(strings.TrimSpace(repo.Format))
	repo.Type = strings.ToLower(strings.TrimSpace(repo.Type))
}
