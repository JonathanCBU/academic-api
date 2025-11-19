package repository

import (
	"academic-api/internal/domain/school"
	"context"

	"github.com/gocraft/dbr/v2"
)

type SchoolRepository struct {
	db *dbr.Session
}

func NewSchoolRepository(db *dbr.Session) *SchoolRepository {
	return &SchoolRepository{db: db}
}

func (r *SchoolRepository) Create(ctx context.Context, s *school.School) error {
	tx, err := r.db.Begin()
	if err != nil {
		return err
	}
	defer tx.RollbackUnlessCommitted()

	err = tx.InsertInto("schools").
		Columns("school_name", "state_code", "district_name").
		Record(s).
		Returning("id", "created_at").
		Load(s)

	if err != nil {
		return err
	}

	return tx.Commit()
}

func (r *SchoolRepository) GetByID(ctx context.Context, id int) (*school.School, error) {
	// Implementation...
}

// Implements school.Repository interface
