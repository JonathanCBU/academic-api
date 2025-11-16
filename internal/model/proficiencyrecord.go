package model

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
	Year     int    `json:"year"`     // Academic year (e.g., 2024)
	State    string `json:"state"`    // State code (e.g., "CA")
	School   string `json:"school"`   // School name (partial match)
	MaxCount int    `json:"maxcount"` // Maximum records to return
}

// GetStatsResponse represents the stats response
type GetStatsResponse struct {
	Records []ProficiencyRecord `json:"records"`
	Count   int                 `json:"count"`
	Total   int                 `json:"total"`
	Filters GetStatsRequest     `json:"filters"`
}

// StatsSummary represents aggregated statistics
type StatsSummary struct {
	TotalSchools           int      `json:"total_schools"`
	TotalRecords           int      `json:"total_records"`
	AverageELAProficiency  *float64 `json:"average_ela_proficiency"`
	AverageMathProficiency *float64 `json:"average_math_proficiency"`
	TotalStudentsTested    int      `json:"total_students_tested"`
}
