package handler

import (
	"academic-api/internal/common"
	"academic-api/internal/service"
	"fmt"
	"net/http"

	"github.com/sirupsen/logrus"
)

type ISchoolReportHandler interface {
	Create(w http.ResponseWriter, r *http.Request)
	Query(w http.ResponseWriter, r *http.Request)
}

type SchoolReportHandler struct {
	ISchoolReportHandler
	service service.ISchoolReportService
}

func NewSchoolReportHandler(service service.ISchoolReportService) *SchoolReportHandler {
	return &SchoolReportHandler{service: service}
}

func (h *SchoolReportHandler) Create(w http.ResponseWriter, r *http.Request) {
	logrus.Info("School Report Handler")
	if r.Body == nil {
		common.WriteBadRequestResponse(w, fmt.Errorf("No body for creating school report object present."))
		return
	}

	schoolObj, err := h.service.Create(r.Body)
	if err != nil {
		common.WriteInternalErrorResponse(w, fmt.Errorf("Failed to create new school report object: %d", err))
		return
	}

	respBody := common.ResponseBody{
		Message: "School report object created.",
		Data:    schoolObj,
	}

	common.WriteCreatedResponse(w, respBody)
}

func (h *SchoolReportHandler) Query(w http.ResponseWriter, r *http.Request) {
	if r.Body == nil {
		common.WriteBadRequestResponse(w, fmt.Errorf("No body for querying school report object present."))
		return
	}

	schools, err := h.service.Query(r.Body)
	if err != nil {
		common.WriteNotFoundResponse(w, fmt.Errorf("Failed to find school report object: %d", err))
		return
	}

	respBody := common.ResponseBody{
		Message: "School report object found.",
		Data:    schools,
	}

	common.WriteOkResponse(w, respBody)
}
