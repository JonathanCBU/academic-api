package api

import "net/http"

type IController interface {
	HandleHealthCheck(w http.ResponseWriter, r *http.Request)
}

type IRouter interface {
	GetRouteHandler() (http.Handler, error)
}
