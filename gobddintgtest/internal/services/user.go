package services

import (
	"context"
	"gobddintgtest/internal/models"
	"gobddintgtest/internal/repos"

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
