package services

import (
	"context"
	"gointegrationtest/internal/models"
	"gointegrationtest/internal/repos"
	"time"

	"github.com/rs/zerolog"
)

type BatchService struct {
	logger zerolog.Logger
	repo   repos.RepoCollection
}

func NewBatchService(logger zerolog.Logger, repo repos.RepoCollection) BatchService {
	return BatchService{
		logger: logger,
		repo:   repo,
	}
}

func (b BatchService) GetBatches(ctx context.Context) ([]models.Batch, error) {
	return b.repo.Batch.GetBatches(ctx)
}

func (b BatchService) CreateBatch(ctx context.Context) (string, error) {
	return b.repo.Batch.InsertBatch(ctx, models.Batch{
		Status:        models.BatchStatusPending,
		DateRequested: time.Now(),
	})
}

func (b BatchService) UpdateBatchStatus(ctx context.Context, id string, status string, errorMessage string) error {
	batch, err := b.repo.Batch.GetBatch(ctx, id)
	if err != nil {
		return err
	}

	batch.Status = status
	batch.ErrorMessage = errorMessage

	return b.repo.Batch.UpdateBatch(ctx, batch)
}
