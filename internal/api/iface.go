package api

import (
	"academic-api/internal/model"
	"io"
	"net/http"
)

type IService interface {
	HealthCheck(reqBody io.ReadCloser) (map[string]any, error)
	GetSchools(req model.GetSchoolsRequest) (*model.GetSchoolsResponse, error)
	GetStats(req model.GetStatsRequest) (*model.GetStatsResponse, error)
}

type IController interface {
	HandleHealthCheck(w http.ResponseWriter, r *http.Request)
	HandleGetSchools(w http.ResponseWriter, r *http.Request)
	HandleGetStats(w http.ResponseWriter, r *http.Request)
}

type IRouter interface {
	GetRouteHandler() (http.Handler, error)
}
