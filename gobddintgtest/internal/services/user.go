package services

import (
	"context"
	"gobddintgtest/internal/models"
	"gobddintgtest/internal/repos"

	"github.com/rs/zerolog"
)

type UserService struct {
	logger   zerolog.Logger
	userRepo repos.UserRepo
}

func NewUserService(logger zerolog.Logger, userRepo repos.UserRepo) UserService {
	return UserService{
		logger:   logger,
		userRepo: userRepo,
	}
}

func (u UserService) GetUsers(ctx context.Context) ([]models.User, error) {
	return u.userRepo.GetUsers()
}
