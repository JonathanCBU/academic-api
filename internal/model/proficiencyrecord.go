package model

import (
	"database/sql"
	"fmt"
	"strings"
)

// ProficiencyRecord represents a single proficiency data point
type ProficiencyRecord struct {
	DataID           int      `json:"data_id"`
	SchoolID         int      `json:"school_id"`
	SchoolName       string   `json:"school_name"`
	StateCode        string   `json:"state_code"`
	DistrictName     string   `json:"district_name"`
	AcademicYear     int      `json:"academic_year"`
	Subject          string   `json:"subject"`
	GradeLevel       string   `json:"grade_level"`
	DemographicGroup string   `json:"demographic_group"`
	NTested          *int     `json:"n_tested"`
	NProficient      *int     `json:"n_proficient"`
	PctProficient    *float64 `json:"pct_proficient"`
	CollectedAt      string   `json:"collected_at"`
}

// GetStatsRequest represents the request body for getting stats
type GetStatsRequest struct {
	Year   int    `json:"year"`   // Academic year (e.g., 2024)
	State  string `json:"state"`  // State code (e.g., "CA")
	School string `json:"school"` // School name (partial match)
}

// GetStatsResponse represents the stats response
type GetStatsResponse struct {
	Records []ProficiencyRecord `json:"records"`
	Count   int                 `json:"count"`
	Filters GetStatsRequest     `json:"filters"`
}

func (req *GetStatsRequest) Validate() error {
	// TODO: implement validation for stats requests
	return nil
}

// TODO: implement pagination
func (req *GetStatsRequest) Query(db *sql.DB) (*GetStatsResponse, error) {
	// Build query with joins
	query := `
		SELECT 
			p.data_id,
			p.school_id,
			s.school_name,
			s.state_code,
			s.district_name,
			p.academic_year,
			p.subject,
			p.grade_level,
			p.demographic_group,
			p.n_tested,
			p.n_proficient,
			p.pct_proficient,
			p.collected_at
		FROM proficiency_data p
		INNER JOIN schools s ON p.school_id = s.school_id
	`

	var whereClauses []string
	var args []interface{}

	// Add filters
	if req.Year > 0 {
		whereClauses = append(whereClauses, "p.academic_year = ?")
		args = append(args, req.Year)
	}

	if req.State != "" {
		whereClauses = append(whereClauses, "s.state_code = ?")
		args = append(args, req.State)
	}

	if req.School != "" {
		whereClauses = append(whereClauses, "s.school_name LIKE ?")
		args = append(args, "%"+req.School+"%")
	}

	// Add WHERE clause if filters exist
	if len(whereClauses) > 0 {
		query += " WHERE " + strings.Join(whereClauses, " AND ")
	}

	// Add ORDER BY
	query += " ORDER BY s.school_id"

	// Execute query
	rows, err := db.Query(query, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to query proficiency data: %w", err)
	}
	defer rows.Close()

	// Scan results
	var records []ProficiencyRecord
	for rows.Next() {
		var record ProficiencyRecord
		err := rows.Scan(
			&record.DataID,
			&record.SchoolID,
			&record.SchoolName,
			&record.StateCode,
			&record.DistrictName,
			&record.AcademicYear,
			&record.Subject,
			&record.GradeLevel,
			&record.DemographicGroup,
			&record.NTested,
			&record.NProficient,
			&record.PctProficient,
			&record.CollectedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan proficiency row: %w", err)
		}
		records = append(records, record)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating proficiency rows: %w", err)
	}

	response := &GetStatsResponse{
		Records: records,
		Count:   len(records),
		Filters: *req,
	}

	return response, nil

}
