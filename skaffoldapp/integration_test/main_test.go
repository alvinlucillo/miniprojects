package integration_test

import (
	"context"
	"log"
	"os"
	"skaffoldapp/integration_test/utils"
	"skaffoldapp/internal/controllers"
	"skaffoldapp/internal/repos"
	"skaffoldapp/internal/services"
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
	// testIntgration := os.Getenv("TEST_INTEGRATION")
	// // Use with TEST_INTEGRATION=true to run integration tests
	// if testIntgration != "true" {
	// 	log.Println("Skipping integration tests")
	// 	os.Exit(0)
	// }
	mongoURI := utils.SetupMongoDB()

	var err error
	utils.MongoClient, err = mongo.Connect(context.TODO(), options.Client().ApplyURI(mongoURI))
	if err != nil {
		log.Fatalf("failed to connect to MongoDB: %v", err)
	}

	// Start Azurite container
	azuriteContainer, blobEndpoint, err := utils.SetupAzuriteContainer()
	if err != nil {
		log.Fatalf("failed to setup Azurite: %v", err)
	}

	// Set environment variables for Azurite compatibility
	os.Setenv("AZURE_STORAGE_ACCOUNT", utils.DefaultAzureAccountName)
	os.Setenv("AZURE_STORAGE_KEY", utils.DefaultAzureBlobKey)
	os.Setenv("AZURE_STORAGE_BLOB_ENDPOINT", blobEndpoint)
	os.Setenv("AZURE_STORAGE_CONTAINER_NAME", utils.TestContainerName)

	// Initialize Azure Manager with Azurite
	azureManager, err := services.NewAzureManager()
	if err != nil {
		log.Fatalf("failed to create AzureManager: %v", err)
	}

	logger := zerolog.New(os.Stdout).With().Timestamp().Logger()

	repoCollection = repos.NewRepoCollection(utils.MongoClient)
	userService = services.NewUserService(logger, repoCollection)
	userController = controllers.NewUsersController(userService)
	batchService = services.NewBatchService(logger, repoCollection, azureManager)

	// Run tests
	code := m.Run()

	// Cleanup
	utils.MongoClient.Disconnect(context.TODO())
	log.Println("Stopping MongoDB test container...")
	utils.TerminateMongoDB()

	_ = azuriteContainer.Terminate(context.Background())

	os.Exit(code)
}
