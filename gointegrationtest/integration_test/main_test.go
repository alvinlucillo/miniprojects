package integration_test

import (
	"context"
	"gointegrationtest/integration_test/utils"
	"gointegrationtest/internal/controllers"
	"gointegrationtest/internal/repos"
	"gointegrationtest/internal/services"
	"log"
	"os"
	"testing"

	"github.com/rs/zerolog"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var (
	repoCollection repos.RepoCollection
	userService    services.UserService
	userController controllers.UsersController
	batchService   services.BatchService
)

func TestMain(m *testing.M) {
	// Start MongoDB test container
	mongoURI := utils.SetupMongoDB()

	// Initialize MongoDB client
	var err error
	utils.MongoClient, err = mongo.Connect(context.TODO(), options.Client().ApplyURI(mongoURI))
	if err != nil {
		log.Fatalf("Failed to connect to MongoDB: %v", err)
	}

	// Start Azurite container
	azuriteContainer, blobEndpoint, err := utils.SetupAzuriteContainer()
	if err != nil {
		log.Fatalf("Failed to setup Azurite: %v", err)
	}

	// Set environment variables for Azurite compatibility
	os.Setenv("AZURE_STORAGE_ACCOUNT", "devstoreaccount1")
	os.Setenv("AZURE_STORAGE_KEY", "Eby8vdM02xNOcqFlqUwJPLlmEtlCDlQvFvhH6c3aqw4=")
	os.Setenv("AZURE_STORAGE_BLOB_ENDPOINT", blobEndpoint)

	// Initialize Azure Manager with Azurite
	azureManager, err := services.NewAzureManager()
	if err != nil {
		log.Fatalf("Failed to create AzureManager: %v", err)
	}

	// Initialize shared repository collection
	repoCollection = repos.NewRepoCollection(utils.MongoClient)

	logger := zerolog.New(os.Stdout).With().Timestamp().Logger()

	// Initialize shared service
	userService = services.NewUserService(logger, repoCollection)

	// Initialize shared controller
	userController = controllers.NewUsersController(userService)

	batchService = services.NewBatchService(logger, repoCollection, azureManager)

	// Run tests
	code := m.Run()

	// Cleanup
	utils.MongoClient.Disconnect(context.TODO())

	// Cleanup MongoDB after all tests
	log.Println("Stopping MongoDB test container...")
	utils.TerminateMongoDB()

	_ = azuriteContainer.Terminate(context.Background()) // Ensure container cleanup

	// Exit with the test result code
	os.Exit(code)
}
