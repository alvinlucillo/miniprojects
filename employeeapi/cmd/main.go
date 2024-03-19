package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"employeeapi/internal/handlers"
	"employeeapi/internal/services"

	"github.com/gin-gonic/gin"
	"github.com/kelseyhightower/envconfig"
	"github.com/rs/zerolog"
)

type Env struct {
	ServerPort string `envconfig:"SERVER_PORT" required:"true" default:"9000"`
}

func main() {
	// Set up logger
	logger := zerolog.New(os.Stdout).With().Timestamp().Logger()
	l := logger.With().Str("package", "main").Logger()

	// Process env vars
	var cfg Env
	err := envconfig.Process("", &cfg)
	if err != nil {
		l.Fatal().Err(err).Msg("failed to process env vars")
	}

	// Set up service and handler
	svc := services.NewService(logger, true)
	h := handlers.NewHandler(logger, svc)

	// Set up server and routes
	r := gin.Default()
	h.SetupRoutes(r)

	srv := &http.Server{
		Addr:    ":" + cfg.ServerPort,
		Handler: r,
	}

	// Start server
	go func() {
		l.Info().Msg("Serving at localhost:" + cfg.ServerPort)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			l.Fatal().Err(err)
		}
	}()

	// Wait for interrupt signal to gracefully shutdown the server
	quit := make(chan os.Signal, 1)

	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	<-quit

	l.Info().Msg("Shutting down server...")

	// Gracefully shutdown the server with a timeout of 5 seconds
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		l.Fatal().Err(err).Msg("Server forced to shutdown")
	}

	l.Info().Msg("Server exited")
}
