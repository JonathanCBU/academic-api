package school

import (
	"academic-api/internal/domain"
	"fmt"

	"github.com/gocraft/dbr/v2"
)

type SchoolRequest struct {
	domain.Request
	StateCode    *string `json:"state_code"`
	DistrictName *string `json:"district_name"`
}

type SchoolResponse struct {
	domain.ApiResponse
	Data []*School
}

func (r *SchoolRequest) ValidateFilter() error {
	if r.StateCode != nil && len(*r.StateCode) != 2 {
		return fmt.Errorf("State code not valid.")
	}

	return nil
}

func (r *SchoolRequest) ApplyFilters(query *dbr.SelectStmt) *dbr.SelectStmt {
	if r.Id != nil {
		query = query.Where("id = ?", *r.Id)
	}

	if r.StateCode != nil {
		query = query.Where("state_code = ?", *r.StateCode)
	}

	if r.DistrictName != nil {
		query = query.Where("district_name = ?", *r.DistrictName)
	}

	return query
}

func (r *SchoolRequest) ApplyCursors(query *dbr.SelectStmt, response *SchoolResponse) (*dbr.SelectStmt, *SchoolResponse) {
	return domain.ApplyCursors(&r.Request, query, response, func(resp *SchoolResponse) *domain.ApiResponse {
		return &resp.ApiResponse
	})
}

func (r *SchoolRequest) Query(db *dbr.Tx) (*SchoolResponse, error) {
	return domain.Query(
		r,
		db,
		"school",
		func() *SchoolResponse { return &SchoolResponse{} },
		func(req *SchoolRequest) *domain.Request { return &req.Request },
		func(resp *SchoolResponse) *domain.ApiResponse { return &resp.ApiResponse },
		func(resp *SchoolResponse) interface{} { return &resp.Data },
		r.ValidateFilter,
		r.ApplyFilters,
	)
}
