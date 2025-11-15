package api

import "net/http"

type IController interface {
	HandleHealthCheck(w http.ResponseWriter, r *http.Request)
}
