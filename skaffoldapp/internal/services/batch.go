package services

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"skaffoldapp/internal/database"
	"skaffoldapp/internal/models"
	"skaffoldapp/internal/repos"
	"time"

	"github.com/rs/zerolog"
)

type BatchService struct {
	logger       zerolog.Logger
	repo         repos.RepoCollection
	azureManager AzureManager
}

func NewBatchService(logger zerolog.Logger, repo repos.RepoCollection, azureManager AzureManager) BatchService {
	return BatchService{
		logger:       logger,
		repo:         repo,
		azureManager: azureManager,
	}
}

func (b BatchService) GetGenerateDBExportRequests(ctx context.Context) ([]models.DBExport, error) {
	return b.repo.Batch.GetDBExports(ctx)
}

func (b BatchService) GenerateDBExport(ctx context.Context) (models.DBExport, error) {
	dbExport := models.DBExport{
		Status:        models.BatchStatusPending,
		DateRequested: time.Now().Truncate(time.Millisecond),
	}

	id, err := b.repo.Batch.InsertDBExports(ctx, dbExport)
	if err != nil {
		return dbExport, err
	}

	dbExport.ID = id

	users, err := b.repo.User.GetUsers(ctx)
	if err != nil {
		return dbExport, fmt.Errorf("failed to get users: %w", err)
	}

	if len(users) == 0 {
		dbExport.Status = models.BatchStatusError
		dbExport.ErrorMessage = "No users found"
		if err := b.repo.Batch.UpdateDBExport(ctx, dbExport); err != nil {
			return dbExport, fmt.Errorf("failed to update dbExport: %w", err)
		}
		return dbExport, nil
	}

	// Generate sqlitedb
	dbExportName := time.Now().Format("20060102150405") + ".db"
	filepath.Join(os.TempDir(), dbExportName)
	defer os.Remove(dbExportName)

	if err := database.CreateUsersDB(users, dbExportName); err != nil {
		dbExport.Status = models.BatchStatusError
		dbExport.ErrorMessage = "Failed to create db: " + err.Error()
		if err := b.repo.Batch.UpdateDBExport(ctx, dbExport); err != nil {
			return dbExport, fmt.Errorf("updating dbExport: %w", err)
		}
		return dbExport, fmt.Errorf("creating db: %w", err)
	}

	if err := b.azureManager.UploadFile(ctx, dbExportName, dbExportName); err != nil {
		dbExport.Status = models.BatchStatusError
		dbExport.ErrorMessage = "Failed to upload file: " + err.Error()
		if err := b.repo.Batch.UpdateDBExport(ctx, dbExport); err != nil {
			return dbExport, fmt.Errorf("updating dbExport: %w", err)
		}
		return dbExport, fmt.Errorf("uploading file: %w", err)
	}
	dbExport.FileName = dbExportName

	return dbExport, nil
}
