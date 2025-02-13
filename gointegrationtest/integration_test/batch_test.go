package integration_test

import (
	"context"
	"gointegrationtest/integration_test/utils"
	"gointegrationtest/internal/models"
	"testing"
	"time"

	"github.com/brianvoe/gofakeit/v7"
	"github.com/stretchr/testify/require"
)

func TestGetBatches(t *testing.T) {
	// Given
	var batches []interface{}
	batches = append(batches, models.Batch{
		DateRequested: gofakeit.Date().Truncate(time.Millisecond),
		Status:        models.BatchStatusPending,
	})
	batches = append(batches, models.Batch{
		DateRequested: gofakeit.Date().Truncate(time.Millisecond),
		Status:        models.BatchStatusError,
		ErrorMessage:  gofakeit.Sentence(10),
	})
	if err := utils.InsertBatches(context.TODO(), batches); err != nil {
		require.NoError(t, err)
	}

	// When
	result, err := batchService.GetBatches(context.TODO())

	// Then
	require.NoError(t, err)
	require.NotEmpty(t, batches)
	require.Equal(t, len(batches), len(result))
	for i := range batches {
		require.Equal(t, batches[i].(models.Batch).DateRequested.Truncate(time.Millisecond), result[i].DateRequested.Truncate(time.Millisecond))
		require.Equal(t, batches[i].(models.Batch).Status, result[i].Status)
		require.Equal(t, batches[i].(models.Batch).ErrorMessage, result[i].ErrorMessage)
	}

	err = utils.CleanupMongoDB()
	if err != nil {
		t.Fatalf("Failed to clean up MongoDB: %v", err)
	}
}
