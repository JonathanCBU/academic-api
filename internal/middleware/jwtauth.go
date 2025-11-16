package middleware

import (
	"academic-api/internal/common"
	"fmt"
	"net/http"
	"strings"

	"github.com/sirupsen/logrus"
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
			token, err := mw.extractAuth(req.Header)
			if err != nil {
				logrus.WithError(err).Error("Failed to extract JWT.")
				common.WriteForbiddenResponse(w, err)
				return
			}

			tokenValid, err := mw.CheckAuth(token)
			if err != nil {
				logrus.WithError(err).Error("Failed to parse JWT.")
				common.WriteUnauthorizedResponse(w, err)
				return
			}
			if !tokenValid {
				logrus.Error("Login attempted with invalidate JWT.")
				common.WriteUnauthorizedResponse(w, err)
				return
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

func (mw *JwtMiddleware) extractAuth(header http.Header) (string, error) {
	token := header.Get(mw.tokenHeaderName)
	if token == "" {
		return "", fmt.Errorf("Auth header not present.")
	}
	if !strings.HasPrefix(token, mw.tokenPrefix) {
		return "", fmt.Errorf("Malformed auth header.")
	}

	token = strings.TrimPrefix(token, mw.tokenPrefix)
	token = strings.TrimSpace(token)
	return token, nil
}
