package dto

import (
	"academic-api/internal/domain/school"
	"time"
)

type CreateSchoolRequest struct {
	SchoolName   string `json:"school_name" validate:"required"`
	StateCode    string `json:"state_code" validate:"required,len=2"`
	DistrictName string `json:"district_name"`
}

func (r *CreateSchoolRequest) ToDomain() *school.School {
	return &school.School{
		SchoolName:   r.SchoolName,
		StateCode:    r.StateCode,
		DistrictName: r.DistrictName,
	}
}

type SchoolResponse struct {
	ID           int    `json:"id"`
	SchoolName   string `json:"school_name"`
	StateCode    string `json:"state_code"`
	DistrictName string `json:"district_name"`
	CreatedAt    string `json:"created_at"`
}

func FromSchool(s *school.School) SchoolResponse {
	return SchoolResponse{
		ID:           s.ID,
		SchoolName:   s.SchoolName,
		StateCode:    s.StateCode,
		DistrictName: s.DistrictName,
		CreatedAt:    s.CreatedAt.Format(time.RFC3339),
	}
}
