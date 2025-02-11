package utils

import (
	"context"
	"gobddintgtest/internal/database"
)

func InsertUsers(ctx context.Context, users []interface{}) error {
	_, err := MongoClient.Database(database.DB_NAME).Collection(database.USER_COLLECTION).InsertMany(ctx, users)
	return err
}
