package handler

import (
	"academic-api/internal/middleware"
	"net/http"
)

type IRouter interface {
	GetRouteHandler() (http.Handler, error)
}

type Router struct {
	scraperController IScraperController
	schoolController  ISchoolHandler
	auth              middleware.IAuthMiddleware
}

func NewRouter(scraperController IScraperController, schoolController ISchoolHandler, auth middleware.IAuthMiddleware) *Router {
	return &Router{
		scraperController: scraperController,
		schoolController:  schoolController,
		auth:              auth,
	}
}
