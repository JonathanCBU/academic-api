package middleware

import (
	"fmt"
	"net/http"
	"strings"
)

const (
	minTokenLen int = 5
)

type JwtMiddleware struct {
	tokenHeaderName string
	tokenPrefix     string
}

func NewJwtMiddleware(tokenHeaderName string, tokenPrefix string) *JwtMiddleware {
	return &JwtMiddleware{
		tokenHeaderName: tokenHeaderName,
		tokenPrefix:     tokenPrefix,
	}
}

func (mw *JwtMiddleware) GetMiddleware() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
			token := req.Header.Get(mw.tokenHeaderName)
			token = strings.TrimPrefix(token, mw.tokenPrefix)

			tokenValid, err := mw.CheckAuth(token)
			if err != nil {
				// TODO: write error response
			}
			if !tokenValid {
				// TODO : write invalid auth response
			}

			next.ServeHTTP(w, req)
		})
	}
}

func (mw *JwtMiddleware) CheckAuth(token string) (bool, error) {
	// TODO: implement JWT auth. This is just placeholder for now
	if len(token) < minTokenLen {
		return false, fmt.Errorf("Token min length %d", minTokenLen)
	}
	return true, nil
}
