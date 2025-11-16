package model

import "fmt"

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
	State    string `json:"state"`    // State abbreviation (e.g., "CA", "TX")
	MaxCount int    `json:"maxcount"` // Maximum number of schools to return
}

// GetSchoolsResponse represents the paginated response
type GetSchoolsResponse struct {
	Schools []School `json:"schools"`
	Count   int      `json:"count"`
	Total   int      `json:"total"`
}

func (req *GetSchoolsRequest) Validate() error {
	// State code format
	if req.State != "" && len(req.State) != 2 {
		return fmt.Errorf("Invalid state abbreviation.")
	}
	return nil
}
