package school

import "context"

type Service struct {
	repo Repository
}

func NewService(repo Repository) *Service {
	return &Service{repo: repo}
}

func (s *Service) CreateSchool(ctx context.Context, school *School) error {
	// Business validation
	if err := school.Validate(); err != nil {
		return err
	}

	// Additional business rules
	// e.g., check for duplicates, apply policies

	return s.repo.Create(ctx, school)
}

func (s *Service) GetSchool(ctx context.Context, id int) (*School, error) {
	return s.repo.GetByID(ctx, id)
}
