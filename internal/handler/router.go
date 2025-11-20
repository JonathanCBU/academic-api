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
	schoolHandler ISchoolHandler
	auth          middleware.IAuthMiddleware
}

func NewRouter(schoolHandler ISchoolHandler, auth middleware.IAuthMiddleware) *Router {
	return &Router{
		schoolHandler: schoolHandler,
		auth:          auth,
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

	router.Use(r.auth.GetMiddleware())

	return router, nil
}
