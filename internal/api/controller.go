package api

import (
	"academic-api/internal/logger"
	"net/http"
)

type Controller struct {
	logCtx *logger.LoggingContext
	client *http.Client
}
