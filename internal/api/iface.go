package api

import (
	"academic-api/internal/model"
	"io"
	"net/http"
)

type IService interface {
	HealthCheck(reqBody io.ReadCloser) (map[string]any, error)
	GetSchools(reader *model.SchoolReader) ([]model.School, error)
	CreateSchool(school *model.School) error
}

type IController interface {
	HandleHealthCheck(w http.ResponseWriter, r *http.Request)
	HandleGetSchools(w http.ResponseWriter, r *http.Request)
	HandleCreateSchool(w http.ResponseWriter, r *http.Request)
}

type IRouter interface {
	GetRouteHandler() (http.Handler, error)
}
