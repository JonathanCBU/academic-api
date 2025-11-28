package schoolreport

import (
	"academic-api/internal/domain"

	"github.com/gocraft/dbr/v2"
)

type SchoolReportRequest struct {
	domain.Request
	SchoolId         *int    `json:"school_id"`
	AcademicYear     *int    `json:"academic_year"`
	Subject          *string `json:"subject"`
	GradeLevel       *string `json:"grade_level"`
	DemographicGroup *string `json:"demographic_group"`
}

type SchoolReportResponse struct {
	domain.ApiResponse
	Data []*SchoolReport
}

func (r *SchoolReportRequest) ValidateFilter() error {
	// TODO: implement filtering validation based on constant lists like in ValidateCreate
	return nil
}

func (r *SchoolReportRequest) ApplyFilters(query *dbr.SelectStmt) *dbr.SelectStmt {
	if r.Id != nil {
		query = query.Where("id = ?", *r.Id)
	}

	if r.SchoolId != nil {
		query = query.Where("school_id = ?", *r.SchoolId)
	}

	if r.AcademicYear != nil {
		query = query.Where("academic_year = ?", *r.AcademicYear)
	}

	if r.DemographicGroup != nil {
		query = query.Where("demographic_group = ?", *r.DemographicGroup)
	}

	if r.GradeLevel != nil {
		query = query.Where("grade_level = ?", *r.GradeLevel)
	}

	if r.Subject != nil {
		query.Where("subject = ?", *r.Subject)
	}

	return query
}

func (r *SchoolReportRequest) ApplyCursors(query *dbr.SelectStmt, response *SchoolReportResponse) (*dbr.SelectStmt, *SchoolReportResponse) {
	return domain.ApplyCursors(&r.Request, query, response, func(resp *SchoolReportResponse) *domain.ApiResponse {
		return &resp.ApiResponse
	})
}

func (r *SchoolReportRequest) Query(db *dbr.Tx) (*SchoolReportResponse, error) {
	return domain.Query(
		r,
		db,
		"school",
		func() *SchoolReportResponse { return &SchoolReportResponse{} },
		func(req *SchoolReportRequest) *domain.Request { return &req.Request },
		func(resp *SchoolReportResponse) *domain.ApiResponse { return &resp.ApiResponse },
		func(resp *SchoolReportResponse) interface{} { return &resp.Data },
		r.ValidateFilter,
		r.ApplyFilters,
	)
}
