package api

import (
	"academic-api/internal/middleware"
	"net/http"

	"github.com/gorilla/mux"
)

type Router struct {
	contoller IController
	auth      middleware.IAuthMiddleware
}

func NewRouter(contoller IController, auth middleware.IAuthMiddleware) *Router {
	return &Router{
		contoller: contoller,
		auth:      auth,
	}
}

func (r *Router) GetRouteHandler() (http.Handler, error) {
	router := mux.NewRouter().StrictSlash(true)

	router.
		Path(healthCheckPath).
		Name(healthCheckName).
		Methods(http.MethodPost).
		HandlerFunc(r.contoller.HandleHealthCheck)

	router.
		Path(schoolsPath).
		Name(schoolsPathName).
		Methods(http.MethodPost).
		HandlerFunc(r.contoller.HandleGetSchools)

	router.
		Path(statsPath).
		Name(statsPathName).
		Methods(http.MethodPost).
		HandlerFunc(r.contoller.HandleGetStats)

	router.
		Path(statsSummaryPath).
		Name(statsSummaryPathName).
		Methods(http.MethodPost).
		HandlerFunc(r.contoller.HandleGetStatsSummary)

	router.Use(r.auth.GetMiddleware())

	return router, nil
}
