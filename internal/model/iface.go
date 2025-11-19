package model

import (
	"database/sql"

	"github.com/gocraft/dbr/v2"
)

type Cursors struct {
	Prev int
	Next int
}

type Model struct {
	Id        int          `db:"id"`
	CreatedAt sql.NullTime `db:"created_at"`
	UpdatedAt sql.NullTime `db:"updated_at"`
}

type IModel interface {
	Create(db *dbr.Tx) error
	Validate() error
}

type PaginatedResponse struct {
	models   []IModel
	cursors  Cursors
	pageSize int
}

type IModelReader interface {
	Validate() error
	GetCursors(db *dbr.Tx, id int, pageSize int, next bool) (Cursors, error)
	Query(db *dbr.Tx) ([]IModel, error)
	QueryPage(db *dbr.Tx, cursors Cursors, pageSize int) (PaginatedResponse, error)
}
