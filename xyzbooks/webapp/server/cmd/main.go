package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"alvinlucillo/xyzbooks_webapp/internal/router"
	service "alvinlucillo/xyzbooks_webapp/internal/services"

	"github.com/kelseyhightower/envconfig"
	"github.com/rs/zerolog"
)

type Environment struct {
	DBHost     string `envconfig:"POSTGRES_HOST" required:"true" default:"localhost"`
	DBPort     string `envconfig:"POSTGRES_PORT" required:"true" default:"5432"`
	DBPassword string `envconfig:"POSTGRES_PASSWORD" required:"true" default:"postgres"`
	DBUser     string `envconfig:"POSTGRES_USER" required:"true" default:"postgres"`
	DBName     string `envconfig:"POSTGRES_DB" required:"true" default:"xyzdb"`
	ServerPort string `envconfig:"SERVER_PORT" required:"true" default:"9001"`
}

func main() {
	logger := zerolog.New(os.Stdout).With().Timestamp().Logger()
	l := logger.With().Str("package", "main").Logger()

	var cfg Environment
	err := envconfig.Process("", &cfg)
	if err != nil {
		l.Fatal().Err(err).Msg("failed to process env vars")
	}

	svcConfig := service.ServiceConfig{
		DBHost:     cfg.DBHost,
		DBPort:     cfg.DBPort,
		DBPassword: cfg.DBPassword,
		DBUser:     cfg.DBUser,
		DBName:     cfg.DBName,
		DBType:     service.DBTypePostgres,
		Logger:     logger,
	}

	svc, err := service.NewService(svcConfig)
	if err != nil {
		l.Fatal().Err(err).Msg("failed to create service")
	}

	rt := router.NewRouter(svc, logger)

	srv := &http.Server{
		Handler: rt.GetHTTPHandler(),
		Addr:    fmt.Sprintf(":%s", cfg.ServerPort),
	}

	// Create a channel to listen for OS signals
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)

	// Run the server in a goroutine so that it doesn't block
	go func() {
		l.Printf("Listening on %s", srv.Addr)
		if err := srv.ListenAndServe(); err != http.ErrServerClosed {
			l.Fatal().Err(err).Msg("ListenAndServe(): %s")
		}
	}()

	// Block until we receive our signal
	<-stop

	// Create a deadline to wait for
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Doesn't block if no connections, but will otherwise wait
	// until the timeout deadline
	srv.Shutdown(ctx)

	log.Println("Shutting down gracefully, press Ctrl+C again to force")

}
