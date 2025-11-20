package domain

import "github.com/gocraft/dbr/v2"

type CursorSet struct {
	Prev int `json:"prev"`
	Next int `json:"next"`
}

type IReader interface {
	Validate() error
	Query(db *dbr.Tx) (*IModel, error)
}

type Reader struct {
	Cursors  CursorSet `json:"cursors"`
	PageSize int       `json:"pageSize"`
	Id       int       `json:"id"`
	IReader
}
