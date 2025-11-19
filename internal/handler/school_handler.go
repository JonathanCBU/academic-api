package handler

import (
	"academic-api/internal/domain/school"
	"academic-api/internal/dto"
	"encoding/json"
	"net/http"
)

type ISchoolHandler interface {
	IDbController
}

type SchoolHandler struct {
	service *school.Service
}

func NewSchoolHandler(service *school.Service) *SchoolHandler {
	return &SchoolHandler{service: service}
}

func (h *SchoolHandler) Create(w http.ResponseWriter, r *http.Request) {
	var req dto.CreateSchoolRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, err)
		return
	}

	school := req.ToDomain() // Convert DTO to domain model

	if err := h.service.CreateSchool(r.Context(), school); err != nil {
		writeError(w, http.StatusInternalServerError, err)
		return
	}

	writeJSON(w, http.StatusCreated, dto.FromSchool(school))
}
