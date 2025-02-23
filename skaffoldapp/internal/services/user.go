package services

import (
	"context"
	"fmt"
	"skaffoldapp/internal/models"
	"skaffoldapp/internal/repos"

	"github.com/rs/zerolog"
)

type UserService struct {
	logger zerolog.Logger
	repo   repos.RepoCollection
}

func NewUserService(logger zerolog.Logger, repo repos.RepoCollection) UserService {
	return UserService{
		logger: logger,
		repo:   repo,
	}
}

func (u UserService) GetUsers(ctx context.Context) ([]models.User, error) {
	return u.repo.User.GetUsers(ctx)
}

func (u UserService) CreateUser(ctx context.Context, user models.UserRequest) (string, error) {
	if user.Name == "" {
		return "", fmt.Errorf(models.EmptyFieldErrorFormat, "name")
	}

	return u.repo.User.InsertUser(ctx, models.User{
		Name: user.Name,
	})
}
