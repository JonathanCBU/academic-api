package schoolreport

import (
	"academic-api/internal/domain"
	"database/sql"
	"fmt"
	"slices"
	"strings"
	"time"

	"github.com/gocraft/dbr/v2"
	"github.com/sirupsen/logrus"
)

var validSubjects = []string{"ela", "math"}
var validGradeLevels = []string{"3", "4", "5", "6", "7", "8", "3-8"}
var validDemographicGroups = []string{"all", "black", "hispanic", "economically_disadvantaged"}

type SchoolReport struct {
	domain.Model
	SchoolId         int    `json:"school_id"`
	DataId           int    `json:"data_id"`
	AcademicYear     int    `json:"academic_year"`
	Subject          string `json:"subject"`
	GradeLevel       string `json:"grade_level"`
	DemographicGroup string `json:"demographic_group"`
	NTested          int    `json:"n_tested"`
	NProficient      int    `json:"n_proficient"`
	PctProficient    int    `json:"pct_proficient"`
}

func NewSchoolReport(schoolId int, dataId int, academicYear int, subject string, gradeLevel string, demographicGroup string, nTested int, nProficient int) *SchoolReport {
	pctProficient := nProficient / nTested * 100.00

	return &SchoolReport{
		SchoolId:         schoolId,
		DataId:           dataId,
		AcademicYear:     academicYear,
		Subject:          subject,
		GradeLevel:       gradeLevel,
		DemographicGroup: demographicGroup,
		NTested:          nTested,
		NProficient:      nProficient,
		PctProficient:    pctProficient,
	}
}

func (r *SchoolReport) ValidateCreate() error {
	if !slices.Contains(validSubjects, strings.ToLower(r.Subject)) {
		return fmt.Errorf("Invalid subject for report: %s", r.Subject)
	}

	if !slices.Contains(validGradeLevels, r.GradeLevel) {
		return fmt.Errorf("Invalid grade level for report: %s", r.GradeLevel)
	}

	if !slices.Contains(validDemographicGroups, strings.ToLower(r.DemographicGroup)) {
		return fmt.Errorf("Invalid demographic group: %s", r.DemographicGroup)
	}

	if r.NTested < 0 {
		return fmt.Errorf("Invalid N tested: %d", r.NTested)
	}

	if r.NProficient < 0 {
		return fmt.Errorf("Invalid N proficient: %d", r.NProficient)
	}

	if r.NProficient > r.NTested {
		return fmt.Errorf("N proficient cannot exceed N tested.")
	}

	return nil
}

func (r *SchoolReport) ValidateUpdate() error {
	err := r.ValidateCreate()
	if err != nil {
		return err
	}

	// TODO: validate other updates
	return nil
}

func (r *SchoolReport) Create(db *dbr.Tx) error {
	err := r.ValidateCreate()
	if err != nil {
		logrus.WithError(err).Error("Failed to validate school report for create.")
	}

	// Set timestamps
	now := time.Now()
	r.CreatedAt = domain.NullTime{
		NullTime: sql.NullTime{Time: now, Valid: true},
	}
	r.UpdatedAt = domain.NullTime{
		NullTime: sql.NullTime{Time: now, Valid: true},
	}
	r.IsDeleted = domain.NullBool{
		NullBool: sql.NullBool{Bool: false, Valid: true},
	}

	// Recalculate pct proficient
	r.PctProficient = r.NProficient / r.NTested * 100

	err = db.InsertInto("school_report").
		Columns(
			"school_id",
			"data_id",
			"academic_year",
			"subject",
			"grade_level",
			"demographic_group",
			"n_tested",
			"n_proficient",
			"pct_proficient",
			"created_at",
			"updated_at",
			"is_deleted",
		).
		Record(r).
		Returning(
			"id",
			"created_at",
			"updated_at",
		).Load(r)
	if err != nil {
		logrus.WithError(err).Error("Failed to insert school report to database.")
		return err
	}

	return nil
}

func (r *SchoolReport) Update(db *dbr.Tx) error {
	err := r.ValidateUpdate()
	if err != nil {
		return err
	}

	err = db.Update("school_report").
		Set("school_id", r.SchoolId).
		Set("data_id", r.DataId).
		Set("academic_year", r.AcademicYear).
		Set("subject", r.Subject).
		Set("grade_level", r.GradeLevel).
		Set("demographic_group", r.DemographicGroup).
		Set("n_tested", r.NTested).
		Set("n_proficient", r.NProficient).
		Set("pct_proficient", r.PctProficient).
		Set("updated_at", time.Now()).
		Where("id = ?", r.Id).
		Returning("updated_at").
		Load(r)

	return err
}

func (r *SchoolReport) Delete(db *dbr.Tx) error {
	err := db.Update("school_report").
		Set("is_deleted", true).
		Set("deleted_at", time.Now()).
		Returning("is_deleted", "deleted_at").
		Load(r)

	return err
}
