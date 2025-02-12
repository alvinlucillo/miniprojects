package database

import (
	"context"
	"fmt"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const (
	DB_NAME         = "gobdddb"
	USER_COLLECTION = "users"
)

func NewMongoDatabase() (*mongo.Client, error) {
	// Connect to MongoDB
	clientOptions := options.Client().ApplyURI("mongodb://admin:admin@localhost:27018")
	client, err := mongo.Connect(context.TODO(), clientOptions)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to MongoDB: %w", err)
	}

	// Ping the database to verify connection
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := client.Ping(ctx, nil); err != nil {
		return nil, fmt.Errorf("MongoDB not reachable: %w", err)
	}

	log.Println("Successfully connected to MongoDB")

	return client, nil
}

func FindAll(collection *mongo.Collection, results interface{}) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Query all documents
	cursor, err := collection.Find(ctx, bson.M{})
	if err != nil {
		return err
	}
	defer cursor.Close(ctx)

	// Decode documents into the provided results slice
	if err = cursor.All(ctx, results); err != nil {
		return err
	}

	return nil
}
