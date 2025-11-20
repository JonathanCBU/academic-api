package main

import (
	"academic-api/internal/handler"
	"academic-api/internal/middleware"
	"academic-api/internal/service"
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gocraft/dbr/v2"
	"github.com/sirupsen/logrus"

	"github.com/joho/godotenv"
	_ "github.com/mattn/go-sqlite3"
)

const (
	// Server configuration
	defaultPort            = "8080"
	defaultEnv             = "development"
	defaultLogLevel        = "debug"
	defaultShutdownTimeout = 30 * time.Second
	defaultReadTimeout     = 15 * time.Second
	defaultWriteTimeout    = 15 * time.Second
	defaultIdleTimeout     = 60 * time.Second
)

func getEnv(envVar string, def string) string {
	val := os.Getenv(envVar)
	if val == "" {
		val = def
	}
	return val
}

func init() {
	err := godotenv.Load()
	if err != nil {
		panic("Failed to load environment.")
	}

	env := getEnv("ENV", defaultEnv)
	logLevel := getEnv("LOG_LEVEL", defaultLogLevel)

	logrus.SetFormatter(&logrus.TextFormatter{
		FullTimestamp:   true,
		TimestampFormat: "2006-01-02 15:04:05",
	})

	level, err := logrus.ParseLevel(logLevel)
	if err != nil {
		logrus.Warn("Invalid log level, defaulting to info")
		level = logrus.InfoLevel
	}
	logrus.SetLevel(level)

	// Enable caller reporting in development
	logrus.SetReportCaller(env != "production")
}

func main() {
	// Add logging context
	log := logrus.WithFields(logrus.Fields{
		"service": "academic-api",
	})

	// Connect to DB
	dbPath := getEnv("DB_PATH", "./data/academic_data.db")
	dbConn, err := dbr.Open("sqlite3", dbPath, nil)
	if err != nil {
		log.WithError(err).Fatal("Failed to conect to database.")
	}
	defer dbConn.Close()

	// TODO: introduce logging for SQL ops here
	dbSess := dbConn.NewSession(nil)

	// Init school service and handler
	schoolService := service.NewSchoolService(dbSess)
	schoolHander := handler.NewSchoolHandler(schoolService)

	// Init auth middleware
	jwtMiddleware := middleware.NewJwtMiddleware(middleware.AuthHeaderName, middleware.BearerPrefix)

	// Init router
	router := handler.NewRouter(schoolHander, jwtMiddleware)
	routeHandler, err := router.GetRouteHandler()
	if err != nil {
		log.WithError(err).Fatal("Failed to create router.")
	}

	port := getEnv("PORT", defaultPort)
	srv := &http.Server{
		Addr:         fmt.Sprintf(":%s", port),
		Handler:      routeHandler,
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
		shutdownTimeoutStr := getEnv("SHUTDOWN_TIMEOUT", "30")
		shutdownTimeout, err := time.ParseDuration(shutdownTimeoutStr + "s")
		if err != nil {
			shutdownTimeout = defaultShutdownTimeout
		}
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
