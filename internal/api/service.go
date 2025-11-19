package api

import (
	"academic-api/internal/common"
	"academic-api/internal/model"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/gocraft/dbr/v2"
	"github.com/sirupsen/logrus"
)

type Service struct {
	client *http.Client
	db     *dbr.Session // Use pointer for consistency
}

func NewService(client *http.Client, db *dbr.Session) *Service {
	if client == nil {
		client = &http.Client{}
	}
	return &Service{
		client: client,
		db:     db,
	}
}

func (s *Service) HealthCheck(reqBody io.ReadCloser) (map[string]interface{}, error) {
	// Create outgoing request
	outgoing, err := http.NewRequest("POST", httpBinPostUrl, reqBody)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}
	outgoing.Header.Set("Content-Type", common.ContentTypeHeader)

	// Execute request
	resp, err := s.client.Do(outgoing)
	if err != nil {
		return nil, fmt.Errorf("failed to execute request: %w", err)
	}
	defer resp.Body.Close()

	// Check status code
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	// Read response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	// Parse JSON response (httpbin returns JSON)
	var result map[string]interface{}
	if err := json.Unmarshal(body, &result); err != nil {
		// If it's not JSON, return as plain text
		return map[string]interface{}{
			"response": string(body),
			"status":   resp.StatusCode,
		}, nil
	}

	return result, nil
}

// GetSchools retrieves a list of schools with optional filtering
func (s *Service) GetSchools(reader *model.SchoolReader) ([]model.IModel, error) {
	// Validate request
	if err := reader.Validate(); err != nil {
		return nil, fmt.Errorf("invalid request: %w", err)
	}

	// Start transaction
	tx, err := s.db.Begin()
	if err != nil {
		logrus.WithError(err).Error("Failed to start database transaction.")
		return nil, fmt.Errorf("database error: %w", err)
	}
	defer tx.RollbackUnlessCommitted() // Auto-rollback if not committed

	// Query schools
	schools, err := reader.Query(tx)
	if err != nil {
		logrus.WithError(err).Error("Failed to query schools.")
		return nil, fmt.Errorf("query error: %w", err)
	}

	// Commit transaction (read-only, but good practice)
	if err := tx.Commit(); err != nil {
		logrus.WithError(err).Error("Failed to commit transaction.")
		return nil, fmt.Errorf("commit error: %w", err)
	}

	return schools, nil
}

// CreateSchool creates a new school record
func (s *Service) CreateSchool(school *model.School) error {
	// Start transaction
	tx, err := s.db.Begin()
	if err != nil {
		logrus.WithError(err).Error("Failed to start database transaction.")
		return fmt.Errorf("database error: %w", err)
	}
	defer tx.RollbackUnlessCommitted()

	// Create school
	if err := school.Create(tx); err != nil {
		logrus.WithError(err).Error("Failed to create school.")
		return fmt.Errorf("create error: %w", err)
	}

	// Commit transaction
	if err := tx.Commit(); err != nil {
		logrus.WithError(err).Error("Failed to commit transaction.")
		return fmt.Errorf("commit error: %w", err)
	}

	logrus.WithField("school_id", school.Id).Info("School created successfully.")
	return nil
}
