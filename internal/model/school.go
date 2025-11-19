package model

import (
	"fmt"

	"github.com/gocraft/dbr/v2"
)

// School represents a school entity
type School struct {
	Model
	SchoolName   string `db:"school_name" json:"school_name"`
	StateCode    string `db:"state_code" json:"state_code"`
	DistrictName string `db:"district_name" json:"district_name"`
}

func (s *School) Validate() error {
	if s.SchoolName == "" {
		return fmt.Errorf("school_name is required")
	}
	if s.StateCode == "" {
		return fmt.Errorf("state_code is required")
	}
	if len(s.StateCode) != 2 {
		return fmt.Errorf("state_code must be 2 characters")
	}
	return nil
}

func (s *School) Create(db *dbr.Tx) error {
	err := s.Validate()
	if err != nil {
		return err
	}

	err = db.InsertInto("schools").
		Columns("school_name", "state_code", "district_name", "created_at").
		Record(s).
		Returning("id", "created_at").
		Load(s) // Load the ID and created_at back into the struct

	if err != nil {
		return err
	}

	return nil
}

// SchoolReader represents query parameters for filtering schools
type SchoolReader struct {
	State *string `json:"state"` // Nullable - use pointer
}

// Helper function for creating a SchoolReader with state filter
func NewSchoolReaderWithState(state string) *SchoolReader {
	return &SchoolReader{State: &state}
}

// Helper function for creating a SchoolReader with no filters
func NewSchoolReader() *SchoolReader {
	return &SchoolReader{State: nil}
}

func (r *SchoolReader) Validate() error {
	// Only validate if State is provided
	if r.State != nil && len(*r.State) != 2 {
		return fmt.Errorf("invalid state abbreviation: must be 2 characters")
	}
	return nil
}

// Query retrieves schools based on filter criteria
func (r *SchoolReader) Query(db *dbr.Tx) ([]IModel, error) {
	var schools []School

	query := db.Select("*").From("schools")

	// Apply state filter if provided
	if r.State != nil {
		query = query.Where("state_code = ?", *r.State)
	}

	_, err := query.OrderBy("id").Load(&schools)
	if err != nil {
		return nil, err
	}

	// Convert []School to []IModel
	result := make([]IModel, len(schools))
	for i := range schools {
		result[i] = &schools[i]
	}

	return result, nil
}

// GetCursors calculates pagination cursors
func (r *SchoolReader) GetCursors(db *dbr.Tx, id int, pageSize int, next bool) (Cursors, error) {
	// Implementation depends on your pagination strategy
	// This is a placeholder
	return Cursors{}, fmt.Errorf("not implemented")
}

// QueryPage retrieves a paginated set of schools
func (r *SchoolReader) QueryPage(db *dbr.Tx, cursors Cursors, pageSize int) (PaginatedResponse, error) {
	// Implementation depends on your pagination strategy
	// This is a placeholder
	return PaginatedResponse{}, fmt.Errorf("not implemented")
}
