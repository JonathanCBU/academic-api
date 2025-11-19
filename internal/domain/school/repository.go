package school

import "context"

// Repository defines the interface for school data access
type Repository interface {
	Create(ctx context.Context, school *School) error
	GetByID(ctx context.Context, id int) (*School, error)
	List(ctx context.Context, filter ListFilter) ([]*School, error)
	Update(ctx context.Context, school *School) error
	Delete(ctx context.Context, id int) error
}

type ListFilter struct {
	StateCode *string
	Limit     int
	Offset    int
}
