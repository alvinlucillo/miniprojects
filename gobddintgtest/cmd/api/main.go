package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"gobddintgtest/internal/controllers"
	"gobddintgtest/internal/database"
	"gobddintgtest/internal/repos"
	"gobddintgtest/internal/routers"
	"gobddintgtest/internal/services"

	"github.com/rs/zerolog"
)

func main() {

	logger := zerolog.New(os.Stdout).With().Timestamp().Logger()

	mongoDBClient, err := database.NewMongoDatabase()
	if err != nil {
		log.Fatalf("Failed to connect to MongoDB: %v", err)
	}

	repoCollection := repos.NewRepoCollection(mongoDBClient)
	userService := services.NewUserService(logger, repoCollection)
	userController := controllers.NewUsersController(userService)

	mux := http.NewServeMux()
	routers.SetupUserRoutes(mux, userController)

	// Create an HTTP server
	server := &http.Server{
		Addr:    ":4200",
		Handler: mux,
	}

	// Create a channel to listen for OS signals
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)

	// Start the server in a goroutine
	go func() {
		log.Printf("Server running on port: %v", server.Addr)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Server failed: %v", err)
		}
	}()

	// Wait for termination signal
	<-stop
	log.Println("Shutting down server...")

	// Create a timeout context for shutdown
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Gracefully shutdown the server
	if err := server.Shutdown(ctx); err != nil {
		log.Fatalf("Server forced to shutdown: %v", err)
	}

	log.Println("Server exited properly")

}
