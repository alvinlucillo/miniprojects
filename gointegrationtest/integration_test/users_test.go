package integration_test

import (
	"context"
	"gointegrationtest/integration_test/utils"
	"gointegrationtest/internal/models"
	"testing"

	"github.com/brianvoe/gofakeit/v7"
	"github.com/stretchr/testify/require"
)

func TestGetUsers(t *testing.T) {
	// Given
	var users []interface{}
	users = append(users, models.User{
		Name: gofakeit.Name(),
	})
	users = append(users, models.User{
		Name: gofakeit.Name(),
	})
	if err := utils.InsertUsers(context.TODO(), users); err != nil {
		require.NoError(t, err)
	}

	// When
	result, err := userService.GetUsers(context.TODO())

	// Then
	require.NoError(t, err)
	require.NotEmpty(t, users)
	require.Equal(t, len(users), len(result))
	for i := range users {
		require.Equal(t, users[i].(models.User).Name, result[i].Name)
	}

	err = utils.CleanupMongoDB()
	if err != nil {
		t.Fatalf("failed to clean up MongoDB: %v", err)
	}
}

func TestCreateUser(t *testing.T) {
	// Given
	user := models.UserRequest{
		Name: gofakeit.Name(),
	}

	// When
	id, err := userService.CreateUser(context.TODO(), user)

	// Then
	require.NoError(t, err)
	require.NotEmpty(t, id)

	users, err := utils.GetUsers(context.TODO())
	require.NoError(t, err)
	require.NotEmpty(t, users)
	require.Equal(t, user.Name, users[0].Name)
	require.Equal(t, id, users[0].ID.Hex())

	err = utils.CleanupMongoDB()
	if err != nil {
		t.Fatalf("failed to clean up MongoDB: %v", err)
	}
}
