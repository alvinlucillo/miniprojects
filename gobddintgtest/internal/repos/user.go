package repos

import (
	"gobddintgtest/internal/models"

	"go.mongodb.org/mongo-driver/mongo"
)

type UserRepo struct {
	db *mongo.Database
}

func NewUserRepo(db *mongo.Database) UserRepo {
	return UserRepo{
		db: db,
	}
}

func (u UserRepo) GetUsers() ([]models.User, error) {
	return nil, nil
}
