package main

import (
	"academic-api/internal/api"
	"academic-api/internal/middleware"
	"context"
	"database/sql"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/joho/godotenv"
	_ "github.com/mattn/go-sqlite3"
	"github.com/sirupsen/logrus"
)

const (
	// Server configuration
	defaultPort            = "8080"
	defaultShutdownTimeout = 30 * time.Second
	defaultReadTimeout     = 15 * time.Second
	defaultWriteTimeout    = 15 * time.Second
	defaultIdleTimeout     = 60 * time.Second
)

func init() {
	setupLogging()
}

func main() {

	// Setup logging with context
	log := logrus.WithFields(logrus.Fields{
		"service": "academic-api",
	})

	log.Info("Starting Academic Data Collection API...")

	// Load .env
	err := godotenv.Load()
	if err != nil {
		logrus.WithError(err).Fatal("Failed to load environment.")
	}

	// Connect to DB
	dbPath := getEnv("DB_PATH", "./data/academic_data.db")
	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		logrus.WithError(err).Fatal("Failed to open database")
	}
	defer db.Close()

	// Initialize dependencies
	httpClient := createHTTPClient()
	service := api.NewService(httpClient, db)
	controller := api.NewController(service)
	jwtMiddleware := middleware.NewJwtMiddleware(middleware.AuthHeaderName, middleware.BearerPrefix)

	// Setup router
	router := api.NewRouter(controller, jwtMiddleware)
	serverRoutes, err := router.GetRouteHandler()
	if err != nil {
		log.WithError(err).Fatal("Failed to initialize router")
	}

	// Configure server
	port := getEnv("PORT", defaultPort)
	srv := &http.Server{
		Addr:         fmt.Sprintf(":%s", port),
		Handler:      serverRoutes,
		ReadTimeout:  defaultReadTimeout,
		WriteTimeout: defaultWriteTimeout,
		IdleTimeout:  defaultIdleTimeout,
	}

	// Start server in goroutine
	serverErrors := make(chan error, 1)
	go func() {
		log.Infof("Server listening on port %s", port)
		log.Infof("Health check available at http://localhost:%s/health-check", port)

		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			serverErrors <- fmt.Errorf("server failed to start: %w", err)
		}
	}()

	// Setup graceful shutdown
	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, syscall.SIGINT, syscall.SIGTERM)

	// Wait for shutdown signal or server error
	select {
	case err := <-serverErrors:
		log.WithError(err).Fatal("Server error occurred")

	case sig := <-shutdown:
		log.Infof("Received shutdown signal: %v", sig)
		log.Info("Starting graceful shutdown...")

		// Get shutdown timeout from env or use default
		shutdownTimeout := getShutdownTimeout()
		ctx, cancel := context.WithTimeout(context.Background(), shutdownTimeout)
		defer cancel()

		// Attempt graceful shutdown
		if err := srv.Shutdown(ctx); err != nil {
			log.WithError(err).Error("Server shutdown encountered error")

			// Force close if graceful shutdown fails
			if err := srv.Close(); err != nil {
				log.WithError(err).Fatal("Failed to force close server")
			}
			log.Warn("Server forcefully closed")
		} else {
			log.Info("Server shutdown completed successfully")
		}
	}

	log.Info("Academic Data Collection API stopped")
}

// setupLogging configures the logger based on environment
func setupLogging() {
	env := getEnv("ENV", "development")
	logLevel := getEnv("LOG_LEVEL", "debug")

	// Set log format based on environment
	if env == "production" {
		logrus.SetFormatter(&logrus.JSONFormatter{
			TimestampFormat: time.RFC3339,
			FieldMap: logrus.FieldMap{
				logrus.FieldKeyTime:  "timestamp",
				logrus.FieldKeyLevel: "level",
				logrus.FieldKeyMsg:   "message",
			},
		})
	} else {
		logrus.SetFormatter(&logrus.TextFormatter{
			FullTimestamp:   true,
			TimestampFormat: "2006-01-02 15:04:05",
		})
	}

	// Set log level
	level, err := logrus.ParseLevel(logLevel)
	if err != nil {
		logrus.Warn("Invalid log level, defaulting to info")
		level = logrus.InfoLevel
	}
	logrus.SetLevel(level)

	// Enable caller reporting in development
	logrus.SetReportCaller(env != "production")
}

// createHTTPClient creates a configured HTTP client for external requests
func createHTTPClient() *http.Client {
	return &http.Client{
		Timeout: 30 * time.Second,
		Transport: &http.Transport{
			MaxIdleConns:        100,
			MaxIdleConnsPerHost: 10,
			IdleConnTimeout:     90 * time.Second,
		},
	}
}

// getEnv retrieves an environment variable with a fallback default
func getEnv(key, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return fallback
}

// getShutdownTimeout gets the shutdown timeout from environment or returns default
func getShutdownTimeout() time.Duration {
	timeoutStr := getEnv("SHUTDOWN_TIMEOUT", "30")
	timeout, err := time.ParseDuration(timeoutStr + "s")
	if err != nil {
		return defaultShutdownTimeout
	}
	return timeout
}
