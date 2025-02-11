package integration_test

import (
	"context"
	"gobddintgtest/integration_test/utils"
	"gobddintgtest/internal/controllers"
	"gobddintgtest/internal/repos"
	"gobddintgtest/internal/services"
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

	// Initialize shared repository collection
	repoCollection = repos.NewRepoCollection(utils.MongoClient)

	logger := zerolog.New(os.Stdout).With().Timestamp().Logger()

	// Initialize shared service
	userService = services.NewUserService(logger, repoCollection) // Pass your logger if needed

	// Initialize shared controller
	userController = controllers.NewUsersController(userService)

	// Run tests
	code := m.Run()

	// Cleanup
	utils.MongoClient.Disconnect(context.TODO())

	// Cleanup MongoDB after all tests
	log.Println("Stopping MongoDB test container...")
	utils.CleanupMongoDB()

	// Exit with the test result code
	os.Exit(code)
}
