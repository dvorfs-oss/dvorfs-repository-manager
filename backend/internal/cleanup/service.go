package cleanup

import (
	"dvorfs-repository-manager/internal/repository"
	"errors"

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
	if err := s.db.Order("created_at desc").Find(&policies).Error; err != nil {
		return nil, err
	}
	return policies, nil
}

func (s *service) CreateCleanupPolicy(policy *repository.CleanupPolicy) error {
	if policy.ID == uuid.Nil {
		policy.ID = uuid.New()
	}
	return s.db.Create(policy).Error
}

func (s *service) UpdateCleanupPolicy(policy *repository.CleanupPolicy) error {
	result := s.db.Model(&repository.CleanupPolicy{}).Where("id = ?", policy.ID).Updates(map[string]any{
		"name":     policy.Name,
		"criteria": policy.Criteria,
	})
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return ErrCleanupPolicyNotFound
	}
	return nil
}

func (s *service) DeleteCleanupPolicy(policyID string) error {
	result := s.db.Where("id = ?", policyID).Delete(&repository.CleanupPolicy{})
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return ErrCleanupPolicyNotFound
	}
	return nil
}
