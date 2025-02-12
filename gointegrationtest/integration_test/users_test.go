package integration_test

import (
	"context"
	"gointegrationtest/integration_test/utils"
	"gointegrationtest/internal/models"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestGetUsers(t *testing.T) {
	// Given
	var users []interface{}
	users = append(users, models.User{
		Name: "John Doe",
	})
	users = append(users, models.User{
		Name: "Jane Doe",
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
}
