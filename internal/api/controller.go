package api

import (
	"academic-api/internal/common"
	"academic-api/internal/model"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/sirupsen/logrus"
)

type Controller struct {
	service IService
}

func NewController(service IService) *Controller {
	return &Controller{
		service: service,
	}
}

func (c *Controller) HandleHealthCheck(w http.ResponseWriter, r *http.Request) {
	// Call service
	data, err := c.service.HealthCheck(r.Body)
	if err != nil {
		logrus.WithError(err).Error("Health check failed")
		common.WriteBadRequestResponse(w, err)
		return
	}

	// Structure response
	respBody := common.ResponseBody{
		Message: "Successfully hit httpBin.",
		Data:    data,
	}

	// Write response
	common.WriteOkResponse(w, respBody)
}

// HandleGetSchools handles POST /schools/list
func (c *Controller) HandleGetSchools(w http.ResponseWriter, r *http.Request) {
	// Parse request body
	var req model.GetSchoolsRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		logrus.WithError(err).Error("Failed to decode request body.")
		common.WriteBadRequestResponse(w, err)
		return
	}
	defer r.Body.Close()

	// Call service
	response, err := c.service.GetSchools(req)
	if err != nil {
		logrus.WithError(err).Error("Failed to get schools.")
		common.WriteInternalErrorResponse(w, err)
		return
	}

	// Return response
	respBody := common.ResponseBody{
		Message: "Schools retrieved successfully",
		Data:    response,
	}

	common.WriteOkResponse(w, respBody)
}

// HandleGetStats handles POST /stats
func (c *Controller) HandleGetStats(w http.ResponseWriter, r *http.Request) {
	// Parse request body
	var req model.GetStatsRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		logrus.WithError(err).Error("Failed to decode request body")
		common.WriteBadRequestResponse(w, fmt.Errorf("invalid request body: %w", err))
		return
	}
	defer r.Body.Close()

	// Validate inputs
	if req.Year < 0 || req.Year > 2100 {
		logrus.WithField("year", req.Year).Warn("Invalid year")
		common.WriteBadRequestResponse(w, fmt.Errorf("year must be between 0 and 2100"))
		return
	}

	if req.State != "" && len(req.State) != 2 {
		logrus.WithField("state", req.State).Warn("Invalid state code format")
		common.WriteBadRequestResponse(w, fmt.Errorf("state code must be 2 characters (e.g., CA, TX)"))
		return
	}

	// Call service
	response, err := c.service.GetStats(req)
	if err != nil {
		logrus.WithError(err).Error("Failed to get stats")
		common.WriteInternalErrorResponse(w, err)
		return
	}

	// Return response
	respBody := common.ResponseBody{
		Message: "Statistics retrieved successfully",
		Data:    response,
	}

	common.WriteOkResponse(w, respBody)
}
