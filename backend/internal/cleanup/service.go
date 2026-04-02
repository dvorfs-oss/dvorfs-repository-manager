package cleanup

import (
	"errors"
	"strings"

	"dvorfs-repository-manager/internal/repository"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Service interface {
	GetAllCleanupPolicies() ([]repository.CleanupPolicy, error)
	CreateCleanupPolicy(policy *repository.CleanupPolicy) error
	UpdateCleanupPolicy(policy *repository.CleanupPolicy) error
	DeleteCleanupPolicy(policyID string) error
}

type service struct {
	db *gorm.DB
}

func NewService(db *gorm.DB) Service {
	return &service{db: db}
}

var ErrCleanupPolicyNotFound = errors.New("cleanup policy not found")

func (s *service) GetAllCleanupPolicies() ([]repository.CleanupPolicy, error) {
	var policies []repository.CleanupPolicy
	err := s.db.Order("name asc").Find(&policies).Error
	return policies, err
}

func (s *service) CreateCleanupPolicy(policy *repository.CleanupPolicy) error {
	if policy == nil {
		return errors.New("policy is required")
	}
	policy.Name = strings.TrimSpace(policy.Name)
	if policy.Name == "" {
		return errors.New("policy name is required")
	}
	if policy.ID == uuid.Nil {
		policy.ID = uuid.New()
	}
	return s.db.Create(policy).Error
}

func (s *service) UpdateCleanupPolicy(policy *repository.CleanupPolicy) error {
	if policy == nil || policy.ID == uuid.Nil {
		return errors.New("policy id is required")
	}

	var existing repository.CleanupPolicy
	if err := s.db.First(&existing, "id = ?", policy.ID).Error; err != nil {
		return err
	}

	existing.Name = strings.TrimSpace(policy.Name)
	existing.Criteria = policy.Criteria
	return s.db.Save(&existing).Error
}

func (s *service) DeleteCleanupPolicy(policyID string) error {
	id, err := uuid.Parse(strings.TrimSpace(policyID))
	if err != nil {
		return err
	}
	return s.db.Delete(&repository.CleanupPolicy{}, "id = ?", id).Error
}
