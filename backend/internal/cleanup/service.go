package cleanup

import (
	"dvorfs-repository-manager/internal/repository"
	"errors"
	"sync"

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
	mu       sync.RWMutex
	policies map[string]repository.CleanupPolicy
	db       *gorm.DB
}

func NewService(db *gorm.DB) Service {
	return &service{
		policies: make(map[string]repository.CleanupPolicy),
		db:       db,
	}
}

var ErrCleanupPolicyNotFound = errors.New("cleanup policy not found")

func (s *service) GetAllCleanupPolicies() ([]repository.CleanupPolicy, error) {
	if s.db != nil {
		var policies []repository.CleanupPolicy
		if err := s.db.Find(&policies).Error; err != nil {
			return nil, err
		}
		return policies, nil
	}

	s.mu.RLock()
	defer s.mu.RUnlock()
	result := make([]repository.CleanupPolicy, 0, len(s.policies))
	for _, policy := range s.policies {
		result = append(result, policy)
	}
	return result, nil
}

func (s *service) CreateCleanupPolicy(policy *repository.CleanupPolicy) error {
	if s.db != nil {
		policy.ID = uuid.New()
		return s.db.Create(policy).Error
	}

	s.mu.Lock()
	defer s.mu.Unlock()
	if policy.ID == uuid.Nil {
		policy.ID = uuid.New()
	}
	s.policies[policy.ID.String()] = *policy
	return nil
}

func (s *service) UpdateCleanupPolicy(policy *repository.CleanupPolicy) error {
	if s.db != nil {
		var existing repository.CleanupPolicy
		if err := s.db.First(&existing, "id = ?", policy.ID).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return ErrCleanupPolicyNotFound
			}
			return err
		}
		return s.db.Save(policy).Error
	}

	s.mu.Lock()
	defer s.mu.Unlock()
	if _, exists := s.policies[policy.ID.String()]; !exists {
		return ErrCleanupPolicyNotFound
	}
	s.policies[policy.ID.String()] = *policy
	return nil
}

func (s *service) DeleteCleanupPolicy(policyID string) error {
	if s.db != nil {
		result := s.db.Delete(&repository.CleanupPolicy{}, "id = ?", policyID)
		if result.Error != nil {
			return result.Error
		}
		if result.RowsAffected == 0 {
			return ErrCleanupPolicyNotFound
		}
		return nil
	}

	s.mu.Lock()
	defer s.mu.Unlock()
	if _, exists := s.policies[policyID]; !exists {
		return ErrCleanupPolicyNotFound
	}
	delete(s.policies, policyID)
	return nil
}
