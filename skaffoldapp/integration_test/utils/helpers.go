package utils

import (
	"context"
	"fmt"
	"skaffoldapp/internal/database"
	"skaffoldapp/internal/models"

	"go.mongodb.org/mongo-driver/bson"
)

func InsertBatches(ctx context.Context, batches []interface{}) error {
	if _, err := MongoClient.Database(database.DB_NAME).Collection(database.EXPORTED_DB_COLLECTION).InsertMany(ctx, batches); err != nil {
		return fmt.Errorf("inserting batches: %w", err)
	}
	return nil
}

func InsertUsers(ctx context.Context, users []interface{}) error {
	if _, err := MongoClient.Database(database.DB_NAME).Collection(database.USER_COLLECTION).InsertMany(ctx, users); err != nil {
		return fmt.Errorf("inserting users: %w", err)
	}
	return nil
}

func GetUsers(ctx context.Context) ([]models.User, error) {
	var users []models.User
	cursor, err := MongoClient.Database(database.DB_NAME).Collection(database.USER_COLLECTION).Find(ctx, bson.D{})
	if err != nil {
		return nil, fmt.Errorf("getting users: %w", err)
	}
	defer cursor.Close(ctx)

	for cursor.Next(ctx) {
		var user models.User
		if err := cursor.Decode(&user); err != nil {
			return nil, fmt.Errorf("decoding user: %w", err)
		}
		users = append(users, user)
	}

	return users, nil
}
