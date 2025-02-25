package integration_test

import (
	"context"
	"os"
	"path/filepath"
	"skaffoldapp/integration_test/utils"
	"skaffoldapp/internal/database"
	"skaffoldapp/internal/models"
	"skaffoldapp/internal/services"
	"testing"
	"time"

	"github.com/brianvoe/gofakeit/v7"
	"github.com/stretchr/testify/require"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func TestGetBatches(t *testing.T) {
	ctx := context.Background()
	// Given
	var dbExports []interface{}
	dbExports = append(dbExports, models.DBExport{
		DateRequested: gofakeit.Date().Truncate(time.Millisecond),
		Status:        models.BatchStatusPending,
	})
	dbExports = append(dbExports, models.DBExport{
		DateRequested: gofakeit.Date().Truncate(time.Millisecond),
		Status:        models.BatchStatusError,
		ErrorMessage:  gofakeit.Sentence(10),
	})
	if err := utils.InsertBatches(ctx, dbExports); err != nil {
		require.NoError(t, err)
	}

	// When
	result, err := batchService.GetGenerateDBExportRequests(ctx)

	// Then
	require.NoError(t, err)
	require.Equal(t, len(dbExports), len(result))
	for i := range dbExports {
		require.Equal(t, dbExports[i].(models.DBExport).DateRequested.Truncate(time.Millisecond), result[i].DateRequested.Truncate(time.Millisecond))
		require.Equal(t, dbExports[i].(models.DBExport).Status, result[i].Status)
		require.Equal(t, dbExports[i].(models.DBExport).ErrorMessage, result[i].ErrorMessage)
	}

	err = utils.CleanupMongoDB()
	if err != nil {
		t.Fatalf("Failed to clean up MongoDB: %v", err)
	}
}

func TestGenerateDBExport(t *testing.T) {
	ctx := context.Background()
	// Given
	var users []interface{}
	users = append(users, models.User{
		Name: gofakeit.Name(),
	})
	users = append(users, models.User{
		Name: gofakeit.Name(),
	})
	if err := utils.InsertUsers(ctx, users); err != nil {
		require.NoError(t, err)
	}

	// When
	dbExport, err := batchService.GenerateDBExport(ctx)

	// Then
	require.NoError(t, err)

	require.NotEmpty(t, dbExport.FileName, "file name should not be empty")
	require.NotEmpty(t, dbExport.Status, "status should not be empty")
	require.NotEmpty(t, dbExport.DateRequested, "date requested should not be empty")
	require.NotEqual(t, primitive.NilObjectID, dbExport.ID, "id should not be empty")
	require.Empty(t, dbExport.ErrorMessage, "error message should be empty")

	azureManager, err := services.NewAzureManager(utils.DefaultAzureAccountName, utils.DefaultAzureBlobKey, AzBlobEndpoint, utils.TestContainerName)
	require.NoError(t, err, "failed to create AzureManager")

	tmpFilePath := filepath.Join(os.TempDir(), dbExport.FileName)
	defer os.Remove(tmpFilePath)

	err = azureManager.GetBlobFile(ctx, dbExport.FileName, tmpFilePath)
	require.NoError(t, err, "failed to get blob file")

	usersFromAzure, err := database.GetUsersFromDatabase(tmpFilePath)
	require.NoError(t, err, "failed to get users from database")

	require.Equal(t, len(users), len(usersFromAzure), "number of users should be the same")
	for i := range users {
		expectedUser, ok := users[i].(models.User)
		require.True(t, ok, "failed to cast user")
		require.Equal(t, expectedUser.Name, usersFromAzure[i].Name, "user names should be the same")
	}

	err = utils.CleanupMongoDB()
	if err != nil {
		t.Fatalf("failed to clean up MongoDB: %v", err)
	}
}
