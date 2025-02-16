package utils

import (
	"context"
	"fmt"
	"gointegrationtest/internal/database"
	"gointegrationtest/internal/models"

	"go.mongodb.org/mongo-driver/bson"
)

func InsertBatches(ctx context.Context, batches []interface{}) error {
	_, err := MongoClient.Database(database.DB_NAME).Collection(database.EXPORTED_DB_COLLECTION).InsertMany(ctx, batches)
	return err
}

func InsertUsers(ctx context.Context, users []interface{}) error {
	_, err := MongoClient.Database(database.DB_NAME).Collection(database.USER_COLLECTION).InsertMany(ctx, users)
	return err
}

func GetUsers(ctx context.Context) ([]models.User, error) {
	var users []models.User
	cursor, err := MongoClient.Database(database.DB_NAME).Collection(database.USER_COLLECTION).Find(ctx, bson.D{})
	if err != nil {
		return nil, fmt.Errorf("failed to get users: %w", err)
	}
	defer cursor.Close(ctx)

	for cursor.Next(ctx) {
		var user models.User
		if err := cursor.Decode(&user); err != nil {
			return nil, fmt.Errorf("failed to decode user: %w", err)
		}
		users = append(users, user)
	}

	return users, nil
}
