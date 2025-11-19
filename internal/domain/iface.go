package domain

import (
	"database/sql"

	"github.com/gocraft/dbr/v2"
)

type IModel interface {
	// Validates field values of model in memory before creating
	ValidateCreate() error

	// Validates field values of model in memory before updating
	ValidateUpdate() error

	// Creates a new database record from model
	Create(dbr.Tx) error

	// Updates database record with corresponding model primary key
	Update(dbr.Tx) error

	// Soft deletes records in the database
	Delete(dbr.Tx) error
}

type Model struct {
	IModel
	Id        int
	IsDeleted bool
	CreatedAt *sql.NullTime
	UpdatedAt *sql.NullTime
	DeletedAt *sql.NullTime
}
