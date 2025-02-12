package repos

import (
	"context"
	"fmt"
	"gointegrationtest/internal/database"
	"gointegrationtest/internal/models"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type UserRepo struct {
	coll *mongo.Collection
}

func NewUserRepo(db *mongo.Database) UserRepo {
	return UserRepo{
		db.Collection("users"),
	}
}

func (u UserRepo) GetUsers(ctx context.Context) ([]models.User, error) {
	var users []models.User

	if err := database.FindAll(u.coll, &users); err != nil {
		return nil, err
	}
	return users, nil
}

func (u UserRepo) InsertUser(ctx context.Context, user models.User) (string, error) {
	result, err := u.coll.InsertOne(ctx, user)
	if err != nil {
		return "", fmt.Errorf("failed to insert user: %w", err)
	}

	id, ok := result.InsertedID.(primitive.ObjectID)
	if !ok {
		return "", fmt.Errorf("failed to convert ID")
	}

	return id.Hex(), err
}
