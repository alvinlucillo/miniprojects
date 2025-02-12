package repos

import (
	"gointegrationtest/internal/database"

	"go.mongodb.org/mongo-driver/mongo"
)

type RepoCollection struct {
	User UserRepo
}

func NewRepoCollection(client *mongo.Client) RepoCollection {
	db := client.Database(database.DB_NAME)

	return RepoCollection{
		User: NewUserRepo(db),
	}
}
