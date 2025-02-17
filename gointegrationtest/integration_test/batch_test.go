package integration_test

import (
	"context"
	"gointegrationtest/integration_test/utils"
	"gointegrationtest/internal/models"
	"testing"
	"time"

	"github.com/brianvoe/gofakeit/v7"
	"github.com/stretchr/testify/require"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func TestGetBatches(t *testing.T) {
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
	if err := utils.InsertBatches(context.TODO(), dbExports); err != nil {
		require.NoError(t, err)
	}

	// When
	result, err := batchService.GetGenerateDBExportRequests(context.TODO())

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
	dbExport, err := batchService.GenerateDBExport(context.TODO())

	// Then
	require.NoError(t, err)

	require.NotEmpty(t, dbExport.FileName, "file name should not be empty")
	require.NotEmpty(t, dbExport.Status, "status should not be empty")
	require.NotEmpty(t, dbExport.DateRequested, "date requested should not be empty")
	require.NotEqual(t, primitive.NilObjectID, dbExport.ID, "id should not be empty")
	require.Empty(t, dbExport.ErrorMessage, "error message should be empty")

	err = utils.CleanupMongoDB()
	if err != nil {
		t.Fatalf("Failed to clean up MongoDB: %v", err)
	}
}
