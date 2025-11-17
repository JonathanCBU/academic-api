package model

import (
	"database/sql"
	"fmt"
)

// School represents a school entity
type School struct {
	SchoolID     int    `json:"school_id"`
	SchoolName   string `json:"school_name"`
	StateCode    string `json:"state_code"`
	DistrictName string `json:"district_name"`
	NCESID       string `json:"nces_id,omitempty"`
	Address      string `json:"address,omitempty"`
	CreatedAt    string `json:"created_at"`
	UpdatedAt    string `json:"updated_at"`
}

// GetSchoolsRequest represents the request body for listing schools
type GetSchoolsRequest struct {
	State string `json:"state"` // State abbreviation (e.g., "CA", "TX")
}

// GetSchoolsResponse represents the paginated response
type GetSchoolsResponse struct {
	Schools []School `json:"schools"`
	Count   int      `json:"count"`
}

func (req *GetSchoolsRequest) Validate() error {
	// State code format
	if req.State != "" && len(req.State) != 2 {
		return fmt.Errorf("Invalid state abbreviation.")
	}
	return nil
}

// inherit
func (req *GetSchoolsRequest) Query(db *sql.DB) (*GetSchoolsResponse, error) {
	// Build query
	query := `
		SELECT 
			school_id,
			school_name,
			state_code,
			district_name,
			COALESCE(nces_id, '') as nces_id,
			COALESCE(address, '') as address,
			created_at,
			updated_at
		FROM schools
	`

	var args []any
	var whereClause string

	// Add state filter if provided
	if req.State != "" {
		whereClause = " WHERE state_code = ?"
		args = append(args, req.State)
	}

	// Add ORDER BY and LIMIT
	query += whereClause + " ORDER BY school_id"

	// Execute query
	rows, err := db.Query(query, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to query schools: %w", err)
	}
	defer rows.Close()

	// Scan results
	var schools []School
	for rows.Next() {
		var school School
		err := rows.Scan(
			&school.SchoolID,
			&school.SchoolName,
			&school.StateCode,
			&school.DistrictName,
			&school.NCESID,
			&school.Address,
			&school.CreatedAt,
			&school.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan school row: %w", err)
		}
		schools = append(schools, school)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating school rows: %w", err)
	}

	response := &GetSchoolsResponse{
		Schools: schools,
		Count:   len(schools),
	}

	return response, nil

}
