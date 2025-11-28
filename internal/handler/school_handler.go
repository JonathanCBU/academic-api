package handler

import (
	"academic-api/internal/common"
	"academic-api/internal/service"
	"fmt"
	"net/http"
)

type ISchoolHandler interface {
	Create(w http.ResponseWriter, r *http.Request)
	Query(w http.ResponseWriter, r *http.Request)
}

type SchoolHandler struct {
	ISchoolHandler
	service service.ISchoolService
}

func NewSchoolHandler(service service.ISchoolService) *SchoolHandler {
	return &SchoolHandler{service: service}
}

func (h *SchoolHandler) Create(w http.ResponseWriter, r *http.Request) {
	if r.Body == nil {
		common.WriteBadRequestResponse(w, fmt.Errorf("No body for creating school object present."))
		return
	}

	schoolObj, err := h.service.Create(r.Body)
	if err != nil {
		common.WriteInternalErrorResponse(w, fmt.Errorf("Failed to create new school object: %d", err))
		return
	}

	respBody := common.ResponseBody{
		Message: "School object created.",
		Data:    schoolObj,
	}

	common.WriteCreatedResponse(w, respBody)
}

func (h *SchoolHandler) Query(w http.ResponseWriter, r *http.Request) {
	if r.Body == nil {
		common.WriteBadRequestResponse(w, fmt.Errorf("No body for querying school object present."))
		return
	}

	schools, err := h.service.Query(r.Body)
	if err != nil {
		common.WriteNotFoundResponse(w, fmt.Errorf("Failed to find school object: %d", err))
		return
	}

	respBody := common.ResponseBody{
		Message: "School object found.",
		Data:    schools,
	}

	common.WriteOkResponse(w, respBody)
}
