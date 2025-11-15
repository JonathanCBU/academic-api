package middleware

import "net/http"

type IMiddlware interface {
	GetMiddleware() func(http.Handler) http.Handler
}

type IAuthMiddleware interface {
	IMiddlware
	CheckAuth(any) (bool, error)
}
