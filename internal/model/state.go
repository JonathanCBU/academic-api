package model

import (
	"context"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type State struct {
	gorm.Model
	ID   int
	Code string
}

func MigrateState(dbFile string) {
	db, err := gorm.Open(sqlite.Open(dbFile), &gorm.Config{})
	if err != nil {
		panic("Failed to connect database.")
	}

	db.AutoMigrate(&State{})
}

func GenerateStatesDummyData(dbFile string) {
	db, err := gorm.Open(sqlite.Open(dbFile), &gorm.Config{})
	if err != nil {
		panic("Failed to connect database.")
	}
	ctx := context.Background()
	// Create
	err = gorm.G[State](db).Create(ctx, &State{Code: "MA"})
	err = gorm.G[State](db).Create(ctx, &State{Code: "CA"})
	err = gorm.G[State](db).Create(ctx, &State{Code: "TX"})

}
