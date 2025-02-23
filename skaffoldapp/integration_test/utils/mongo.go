package utils

import (
	"context"
	"fmt"
	"log"
	"skaffoldapp/internal/database"

	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

var MongoClient *mongo.Client
var mongoContainer testcontainers.Container

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

func TerminateMongoDB() error {
	log.Println("Stopping MongoDB test container...")
	return mongoContainer.Terminate(context.Background())
}

func CleanupMongoDB() error {
	log.Println("Cleaning up MongoDB...")
	_, _ = MongoClient.Database("test").Collection("users").DeleteMany(context.Background(), nil)

	collections, err := MongoClient.Database(database.DB_NAME).ListCollectionNames(context.TODO(), bson.M{})
	if err != nil {
		return err
	}

	for _, coll := range collections {
		err := MongoClient.Database(database.DB_NAME).Collection(coll).Drop(context.TODO())
		if err != nil {
			return fmt.Errorf("failed to drop collection %s: %w", coll, err)
		}
	}
	return nil
}
