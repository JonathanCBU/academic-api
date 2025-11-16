package api

import (
	"academic-api/internal/common"
	"academic-api/internal/model"
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
)

type Service struct {
	client *http.Client
	db     *sql.DB
}

func NewService(client *http.Client, db *sql.DB) *Service {
	if client == nil {
		client = &http.Client{}
	}
	return &Service{
		client: client,
		db:     db,
	}
}

func (s *Service) HealthCheck(reqBody io.ReadCloser) (map[string]interface{}, error) {
	// Create outgoing request
	outgoing, err := http.NewRequest("POST", httpBinPostUrl, reqBody)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}
	outgoing.Header.Set("Content-Type", common.ContentTypeHeader)

	// Execute request
	resp, err := s.client.Do(outgoing)
	if err != nil {
		return nil, fmt.Errorf("failed to execute request: %w", err)
	}
	defer resp.Body.Close()

	// Check status code
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	// Read response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	// Parse JSON response (httpbin returns JSON)
	var result map[string]interface{}
	if err := json.Unmarshal(body, &result); err != nil {
		// If it's not JSON, return as plain text
		return map[string]interface{}{
			"response": string(body),
			"status":   resp.StatusCode,
		}, nil
	}

	return result, nil
}

// GetSchools retrieves a list of schools with optional filtering
func (s *Service) GetSchools(req model.GetSchoolsRequest) (*model.GetSchoolsResponse, error) {
	// Set default max count
	if req.MaxCount <= 0 {
		req.MaxCount = 50
	}

	// Cap max count to prevent excessive queries
	if req.MaxCount > 100 {
		req.MaxCount = 100
	}

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

	var args []interface{}
	var whereClause string

	// Add state filter if provided
	if req.State != "" {
		whereClause = " WHERE state_code = ?"
		args = append(args, req.State)
	}

	// Add ORDER BY and LIMIT
	query += whereClause + " ORDER BY school_name ASC LIMIT ?"
	args = append(args, req.MaxCount)

	// Execute query
	rows, err := s.db.Query(query, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to query schools: %w", err)
	}
	defer rows.Close()

	// Scan results
	var schools []model.School
	for rows.Next() {
		var school model.School
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

	// Get total count (without limit)
	totalCount, err := s.getTotalSchoolCount(req.State)
	if err != nil {
		return nil, fmt.Errorf("failed to get total count: %w", err)
	}

	response := &model.GetSchoolsResponse{
		Schools: schools,
		Count:   len(schools),
		Total:   totalCount,
	}

	return response, nil
}

// getTotalSchoolCount returns the total number of schools (for pagination info)
func (s *Service) getTotalSchoolCount(stateCode string) (int, error) {
	query := "SELECT COUNT(*) FROM schools"
	var args []interface{}

	if stateCode != "" {
		query += " WHERE state_code = ?"
		args = append(args, stateCode)
	}

	var count int
	err := s.db.QueryRow(query, args...).Scan(&count)
	if err != nil {
		return 0, err
	}

	return count, nil
}

func (s *Service) GetStats(req model.GetStatsRequest) (*model.GetStatsResponse, error) {
	// Set default max count
	if req.MaxCount <= 0 {
		req.MaxCount = 50
	}

	// Cap max count to prevent excessive queries
	if req.MaxCount > 200 {
		req.MaxCount = 200
	}

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

	// Add ORDER BY and LIMIT
	query += " ORDER BY s.school_name, p.subject, p.demographic_group LIMIT ?"
	args = append(args, req.MaxCount)

	// Execute query
	rows, err := s.db.Query(query, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to query proficiency data: %w", err)
	}
	defer rows.Close()

	// Scan results
	var records []model.ProficiencyRecord
	for rows.Next() {
		var record model.ProficiencyRecord
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

	// Get total count (without limit)
	totalCount, err := s.getTotalStatsCount(req)
	if err != nil {
		return nil, fmt.Errorf("failed to get total count: %w", err)
	}

	response := &model.GetStatsResponse{
		Records: records,
		Count:   len(records),
		Total:   totalCount,
		Filters: req,
	}

	return response, nil
}

// getTotalStatsCount returns the total number of records matching filters
func (s *Service) getTotalStatsCount(req model.GetStatsRequest) (int, error) {
	query := `
		SELECT COUNT(*) 
		FROM proficiency_data p
		INNER JOIN schools s ON p.school_id = s.school_id
	`

	var whereClauses []string
	var args []interface{}

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

	if len(whereClauses) > 0 {
		query += " WHERE " + strings.Join(whereClauses, " AND ")
	}

	var count int
	err := s.db.QueryRow(query, args...).Scan(&count)
	if err != nil {
		return 0, err
	}

	return count, nil
}

// GetStatsSummary returns aggregated statistics
func (s *Service) GetStatsSummary(req model.GetStatsRequest) (*model.StatsSummary, error) {
	query := `
		SELECT 
			COUNT(DISTINCT p.school_id) as total_schools,
			COUNT(*) as total_records,
			AVG(CASE WHEN p.subject = 'ela' THEN p.pct_proficient END) as avg_ela,
			AVG(CASE WHEN p.subject = 'math' THEN p.pct_proficient END) as avg_math,
			SUM(COALESCE(p.n_tested, 0)) as total_tested
		FROM proficiency_data p
		INNER JOIN schools s ON p.school_id = s.school_id
	`

	var whereClauses []string
	var args []interface{}

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

	if len(whereClauses) > 0 {
		query += " WHERE " + strings.Join(whereClauses, " AND ")
	}

	var summary model.StatsSummary
	err := s.db.QueryRow(query, args...).Scan(
		&summary.TotalSchools,
		&summary.TotalRecords,
		&summary.AverageELAProficiency,
		&summary.AverageMathProficiency,
		&summary.TotalStudentsTested,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to get stats summary: %w", err)
	}

	return &summary, nil
}
