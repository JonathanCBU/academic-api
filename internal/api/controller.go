package api

import (
	"academic-api/internal/common"
	"academic-api/internal/model"
	"encoding/json"
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
	var reader model.SchoolReader
	if err := json.NewDecoder(r.Body).Decode(&reader); err != nil {
		logrus.WithError(err).Error("Failed to decode request body.")
		common.WriteBadRequestResponse(w, err)
		return
	}
	defer r.Body.Close()

	// Call service - pass pointer
	schools, err := c.service.GetSchools(&reader)
	if err != nil {
		logrus.WithError(err).Error("Failed to get schools.")
		common.WriteInternalErrorResponse(w, err)
		return
	}

	// Return response
	respBody := common.ResponseBody{
		Message: "Schools retrieved successfully",
		Data:    schools,
	}

	common.WriteOkResponse(w, respBody)
}

// HandleCreateSchool handles POST /schools
func (c *Controller) HandleCreateSchool(w http.ResponseWriter, r *http.Request) {
	// Parse request body
	var school model.School
	if err := json.NewDecoder(r.Body).Decode(&school); err != nil {
		logrus.WithError(err).Error("Failed to decode request body.")
		common.WriteBadRequestResponse(w, err)
		return
	}
	defer r.Body.Close()

	// Call service
	if err := c.service.CreateSchool(&school); err != nil {
		logrus.WithError(err).Error("Failed to create school.")
		common.WriteInternalErrorResponse(w, err)
		return
	}

	// Return response with created school (includes generated ID)
	respBody := common.ResponseBody{
		Message: "School created successfully",
		Data:    school,
	}

	w.WriteHeader(http.StatusCreated)
	common.WriteOkResponse(w, respBody)
}
