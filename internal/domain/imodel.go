package domain

import (
	"github.com/gocraft/dbr/v2"
)

type IModel interface {
	// Validates field values of model in memory before creating
	ValidateCreate() error

	// Validates field values of model in memory before updating
	ValidateUpdate() error

	// Creates a new database record from model
	Create(db *dbr.Tx) error

	// Updates database record with corresponding model primary key
	Update(db *dbr.Tx) error

	// Soft deletes records in the database
	Delete(db *dbr.Tx) error
}

type Model struct {
	IModel    `json:"-"`
	Id        int      `json:"id"`
	IsDeleted NullBool `json:"is_deleted"`
	CreatedAt NullTime `json:"created_at"`
	UpdatedAt NullTime `json:"updated_at"`
	DeletedAt NullTime `json:"deleted_at"`
}
