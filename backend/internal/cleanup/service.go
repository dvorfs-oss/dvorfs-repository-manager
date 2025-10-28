package cleanup

import "dvorfs-repository-manager/internal/repository"

type Service interface {
	GetAllCleanupPolicies() ([]repository.CleanupPolicy, error)
	CreateCleanupPolicy(policy *repository.CleanupPolicy) error
	UpdateCleanupPolicy(policy *repository.CleanupPolicy) error
	DeleteCleanupPolicy(policyID string) error
}

type service struct {
	// Add dependencies here, e.g., a cleanup policy repository
}

func NewService() Service {
	return &service{}
}

func (s *service) GetAllCleanupPolicies() ([]repository.CleanupPolicy, error) {
	return []repository.CleanupPolicy{{Name: "test-policy"}}, nil
}

func (s *service) CreateCleanupPolicy(policy *repository.CleanupPolicy) error {
	return nil
}

func (s *service) UpdateCleanupPolicy(policy *repository.CleanupPolicy) error {
	return nil
}

func (s *service) DeleteCleanupPolicy(policyID string) error {
	return nil
}
