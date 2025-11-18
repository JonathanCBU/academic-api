package api

import (
	"academic-api/internal/common"
	"academic-api/internal/model"
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/sirupsen/logrus"
)

type Service struct {
	client *http.Client
	db     *sql.DB
}

func NewService(client *http.Client, db *sql.DB) *Service {
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
func (s *Service) GetSchools(req model.GetSchoolsRequest) (*model.GetSchoolsResponse, error) {

	err := req.Validate()
	if err != nil {
		logrus.WithError(err).Error("Failed to validate schools request body.")
		return nil, err
	}
	response, err := req.Query(s.db)
	if err != nil {
		logrus.WithError(err).Error("Failed to query db for school data")
		return nil, err
	}

	return response, nil
}

func (s *Service) GetStats(req model.GetStatsRequest) (*model.GetStatsResponse, error) {

	err := req.Validate()
	if err != nil {
		logrus.WithError(err).Error("Failed to validate stats request body.")
		return nil, err
	}
	response, err := req.Query(s.db)
	if err != nil {
		logrus.WithError(err).Error("Failed to query db for stats")
	}

	return response, nil
}
