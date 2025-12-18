package domain

import (
	"github.com/gocraft/dbr/v2"
)

type CursorSet struct {
	Prev *int `json:"prev"`
	Next *int `json:"next"`
}

type ApiResponse struct {
	Cursors  CursorSet `json:"cursors"`
	PageSize *int      `json:"page_size"`
}

type IRequest interface {
	ValidateFilter() error
	ApplyFilters(query *dbr.SelectStmt) *dbr.SelectStmt
	ApplyCursors(query *dbr.SelectStmt, response *ApiResponse)
	Query(db *dbr.Tx) (*ApiResponse, error)
}

type Request struct {
	Cursors  CursorSet `json:"cursors"`
	PageSize *int      `json:"page_size"`
	Id       *int      `json:"id"`
	IRequest
}

// Generic cursors application
func ApplyCursors[T any](r *Request, query *dbr.SelectStmt, response *T, getApiResponse func(*T) *ApiResponse) (*dbr.SelectStmt, *T) {
	apiResp := getApiResponse(response)

	if r.PageSize == nil {
		return query, response
	}

	apiResp.PageSize = r.PageSize

	if r.Cursors.Next == nil && r.Cursors.Prev == nil {
		return query, response
	}

	var endId int
	if r.Cursors.Next != nil && *r.Cursors.Next > 0 {
		endId = *r.Cursors.Next + *r.PageSize
		query = query.Where("id >= ?", *r.Cursors.Next).Where("id < ?", endId)
		apiResp.Cursors.Next = &endId
	} else if r.Cursors.Prev != nil && *r.Cursors.Prev > 0 {
		endId = *r.Cursors.Prev - *r.PageSize
		query = query.Where("id > ?", endId).Where("id <= ?", *r.Cursors.Prev)
		apiResp.Cursors.Prev = &endId
	}

	return query, response
}

// Generic Query function
func Query[TReq any, TResp any](
	req TReq,
	db *dbr.Tx,
	tableName string,
	newResponse func() *TResp,
	getRequest func(TReq) *Request,
	getApiResponse func(*TResp) *ApiResponse,
	getDataPtr func(*TResp) interface{},
	validateFilter func() error,
	applyFilters func(*dbr.SelectStmt) *dbr.SelectStmt,
) (*TResp, error) {
	// Validate filters
	err := validateFilter()
	if err != nil {
		return nil, err
	}

	response := newResponse()
	baseReq := getRequest(req)

	// Build query
	query := db.Select("*").From(tableName)
	query = applyFilters(query)
	query, response = ApplyCursors(baseReq, query, response, getApiResponse)

	// Load data
	_, err = query.Load(getDataPtr(response))
	if err != nil {
		return nil, err
	}

	return response, nil
}
