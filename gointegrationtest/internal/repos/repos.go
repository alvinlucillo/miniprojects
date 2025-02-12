package repos

import "go.mongodb.org/mongo-driver/mongo"

type RepoCollection struct {
	User UserRepo
}

func NewRepoCollection(client *mongo.Client) RepoCollection {
	db := client.Database("gobdddb")

	return RepoCollection{
		User: NewUserRepo(db),
	}
}
