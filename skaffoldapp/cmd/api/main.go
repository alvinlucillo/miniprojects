package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"skaffoldapp/internal/controllers"
	"skaffoldapp/internal/database"
	"skaffoldapp/internal/repos"
	"skaffoldapp/internal/routers"
	"skaffoldapp/internal/services"

	"github.com/kelseyhightower/envconfig"
	"github.com/rs/zerolog"
)

type Environment struct {
	MongoDBUsername       string `required:"true" envconfig:"MONGODB_USERNAME"`
	MongoDBPassword       string `required:"true" envconfig:"MONGODB_PASSWORD"`
	MongoDBHost           string `required:"true" envconfig:"MONGODB_HOST"`
	MongoDBPort           string `required:"true" envconfig:"MONGODB_PORT"`
	AzStorageAccount      string `required:"true" envconfig:"AZURE_STORAGE_ACCOUNT"`
	AzStorageKey          string `required:"true" envconfig:"AZURE_STORAGE_KEY"`
	AzStorageBlobEndpoint string `required:"true" envconfig:"AZURE_STORAGE_BLOB_ENDPOINT"`
	AzContainerName       string `required:"true" envconfig:"AZURE_STORAGE_CONTAINER_NAME"`
}

func main() {
	logger := zerolog.New(os.Stdout).With().Timestamp().Logger()

	var env Environment
	if err := envconfig.Process("", &env); err != nil {
		log.Fatalf("failed to load environment variables: %v", err)
	}

	mongoDBClient, err := database.NewMongoDatabase(env.MongoDBHost, env.MongoDBPort, env.MongoDBUsername, env.MongoDBPassword)
	if err != nil {
		log.Fatalf("failed to connect to MongoDB: %v", err)
	}

	azureManager, err := services.NewAzureManager(env.AzStorageAccount, env.AzStorageKey, env.AzStorageBlobEndpoint, env.AzContainerName)
	if err != nil {
		log.Fatalf("failed to create AzureManager: %v", err)
	}

	repoCollection := repos.NewRepoCollection(mongoDBClient)
	userService := services.NewUserService(logger, repoCollection)
	userController := controllers.NewUsersController(userService)
	batchController := controllers.NewBatchController(services.NewBatchService(logger, repoCollection, azureManager))

	mux := http.NewServeMux()
	routers.SetupUserRoutes(mux, userController)
	routers.SetupBatchRoutes(mux, batchController)

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
			log.Printf("Server failed: %v", err)
		}
	}()

	// Wait for termination signal
	<-stop

	// Create a timeout context for shutdown
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	log.Println("Disconnecting database...")
	if err := mongoDBClient.Disconnect(ctx); err != nil {
		log.Printf("Failed to disconnect from MongoDB: %v", err)
	}

	log.Println("Shutting down server...")

	// Gracefully shutdown the server
	if err := server.Shutdown(ctx); err != nil {
		log.Fatalf("Server forced to shutdown: %v", err)
	}

	log.Println("Server exited")
}
