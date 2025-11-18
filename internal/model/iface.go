package model

import "database/sql"

type IResult interface{}

type Cursors struct {
	Prev int
	Next int
}

type PaginatedResult struct {
	cursors Cursors
}

type IRequest interface {
	// Validate incoming request fields
	Validate() error

	// Perform DB query
	Query(db *sql.DB) (IResult, error)
	
	// Perform DB query with cursor pagination
	PaginatedQuery(db *sql.DB, cursor)
}
