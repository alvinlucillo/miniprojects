package repos

import (
	"context"
	"gointegrationtest/internal/database"
	"gointegrationtest/internal/models"

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
