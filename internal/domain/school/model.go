package school

import (
	"academic-api/internal/domain"
	"database/sql"
	"fmt"
	"time"

	"github.com/gocraft/dbr/v2"
	"github.com/sirupsen/logrus"
)

type School struct {
	domain.Model
	SchoolName   string `json:"school_name"`
	StateCode    string `json:"state_code"`
	DistrictName string `json:"district_name"`
}

func NewSchool(name string, state string, district string) *School {
	return &School{
		SchoolName:   name,
		StateCode:    state,
		DistrictName: district,
	}
}

func (s *School) ValidateCreate() error {
	if s.SchoolName == "" {
		return fmt.Errorf("Invalid school name.")
	}
	if len(s.StateCode) != 2 {
		return fmt.Errorf("Invalid state code.")
	}
	return nil
}

func (s *School) ValidateUpdate() error {
	err := s.ValidateCreate()
	if err != nil {
		return err
	}

	// TODO: Prevent update for invalid timestamps or attempts to update deleted records
	return nil
}

func (s *School) Create(db *dbr.Tx) error {
	err := s.ValidateCreate()
	if err != nil {
		logrus.WithError(err).Error("Failed to validate school content for create.")
		return err
	}

	// Set timestamps
	now := time.Now()
	s.CreatedAt = domain.NullTime{
		NullTime: sql.NullTime{Time: now, Valid: true},
	}
	s.UpdatedAt = domain.NullTime{
		NullTime: sql.NullTime{Time: now, Valid: true},
	}
	s.IsDeleted = domain.NullBool{
		NullBool: sql.NullBool{Bool: false, Valid: true},
	}

	err = db.InsertInto("school").
		Columns("school_name", "state_code", "district_name", "created_at", "updated_at", "is_deleted").
		Record(s).
		Returning("id", "created_at", "updated_at").
		Load(s) // Load the ID and created_at back into the struct

	if err != nil {
		logrus.WithError(err).Error("Failed to insert school to database.")
		return err
	}

	return nil

}

func (s *School) Update(db *dbr.Tx) error {
	err := s.ValidateUpdate()
	if err != nil {
		return err
	}

	err = db.Update("school").
		Set("school_name", s.SchoolName).
		Set("state_code", s.StateCode).
		Set("district_name", s.DistrictName).
		Set("updated_at", time.Now()).
		Where("id = ?", s.Id).
		Returning("updated_at").
		Load(s)

	return err
}

func (s *School) Delete(db *dbr.Tx) error {
	err := db.Update("school").
		Set("is_deleted", true).
		Set("deleted_at", time.Now()).
		Returning("is_deleted", "deleted_at").
		Load(s)

	return err
}
