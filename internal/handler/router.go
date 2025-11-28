package handler

import (
	"academic-api/internal/middleware"
	"net/http"

	"github.com/gorilla/mux"
)

type IRouter interface {
	GetRouteHandler() (http.Handler, error)
}

type Router struct {
	schoolHandler       ISchoolHandler
	schoolReportHandler ISchoolReportHandler
	auth                middleware.IAuthMiddleware
}

func NewRouter(schoolHandler ISchoolHandler, schoolReportHandler ISchoolReportHandler, auth middleware.IAuthMiddleware) *Router {
	return &Router{
		schoolHandler:       schoolHandler,
		schoolReportHandler: schoolReportHandler,
		auth:                auth,
	}
}

func (r *Router) GetRouteHandler() (http.Handler, error) {
	router := mux.NewRouter().StrictSlash(true)

	router.
		Path(schoolsPath + "/put").
		Name(schoolsPathName + "Put").
		Methods(http.MethodPost).
		HandlerFunc(r.schoolHandler.Create)

	router.
		Path(schoolsPath + "/get").
		Name(schoolsPathName + "Get").
		Methods(http.MethodPost).
		HandlerFunc(r.schoolHandler.Query)

	router.
		Path(schoolReportsPath + "/put").
		Name(schoolReportsPathName + "Put").
		Methods(http.MethodPost).
		HandlerFunc(r.schoolReportHandler.Create)

	router.
		Path(schoolReportsPath + "/get").
		Name(schoolReportsPathName + "Get").
		Methods(http.MethodPost).
		HandlerFunc(r.schoolReportHandler.Query)

	router.Use(r.auth.GetMiddleware())

	return router, nil
}
