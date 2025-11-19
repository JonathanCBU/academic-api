package main

import (
	"academic-api/internal/domain/school"
	"academic-api/internal/handler"
	"academic-api/internal/repository"
	"academic-api/pkg/database"
	"net/http"
)

func main() {
	// Initialize dependencies
	db := database.Connect()

	// Repository layer
	schoolRepo := repository.NewSchoolRepository(db)

	// Service layer
	schoolService := school.NewService(schoolRepo)

	// Handler layer
	schoolHandler := handler.NewSchoolHandler(schoolService)

	// Setup routes
	router := handler.NewRouter(schoolHandler)

	http.ListenAndServe(":8080", router)
}
