package utils

import (
	"context"
	"fmt"
	"log"

	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
	"go.mongodb.org/mongo-driver/mongo"
)

var MongoClient *mongo.Client
var mongoContainer testcontainers.Container

// SetupMongoDB starts MongoDB and returns the URI
func SetupMongoDB() string {
	ctx := context.Background()
	req := testcontainers.ContainerRequest{
		Image:        "mongo:6.0",
		ExposedPorts: []string{"27017/tcp"},
		WaitingFor:   wait.ForListeningPort("27017/tcp"),
	}

	var err error
	mongoContainer, err = testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: req,
		Started:          true,
	})
	if err != nil {
		log.Fatalf("Failed to start MongoDB container: %v", err)
	}

	// Get MongoDB URI
	host, _ := mongoContainer.Host(ctx)
	port, _ := mongoContainer.MappedPort(ctx, "27017")
	mongoURI := fmt.Sprintf("mongodb://%s:%s", host, port.Port())

	log.Printf("MongoDB Test Container started at %s", mongoURI)

	return mongoURI
}

// CleanupMongoDB stops MongoDB after all tests
func CleanupMongoDB() {
	log.Println("Stopping MongoDB test container...")
	_ = mongoContainer.Terminate(context.Background())
}
