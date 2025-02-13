package repos

import (
	"context"
	"fmt"
	"gointegrationtest/internal/database"
	"gointegrationtest/internal/models"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type BatchRepo struct {
	coll *mongo.Collection
}

func NewBatchRepo(db *mongo.Database) BatchRepo {
	return BatchRepo{
		db.Collection(database.BATCH_COLLECTION),
	}
}

func (b BatchRepo) GetBatches(ctx context.Context) ([]models.Batch, error) {
	var batches []models.Batch

	if err := database.FindAll(b.coll, &batches); err != nil {
		return nil, fmt.Errorf("failed to get batches: %w", err)
	}
	return batches, nil
}

func (b BatchRepo) InsertBatch(ctx context.Context, batch models.Batch) (string, error) {
	result, err := b.coll.InsertOne(ctx, batch)
	if err != nil {
		return "", fmt.Errorf("failed to insert batch: %w", err)
	}

	id := result.InsertedID.(string)
	return id, nil
}

func (b BatchRepo) UpdateBatch(ctx context.Context, batch models.Batch) error {
	_, err := b.coll.UpdateOne(ctx, bson.M{"_id": batch.ID}, bson.M{"$set": batch})
	if err != nil {
		return fmt.Errorf("failed to update batch: %w", err)
	}
	return nil
}

func (b BatchRepo) GetBatch(ctx context.Context, id string) (models.Batch, error) {
	var batch models.Batch

	err := b.coll.FindOne(ctx, bson.M{"_id": id}).Decode(&batch)
	if err != nil {
		return models.Batch{}, fmt.Errorf("failed to get batch: %w", err)
	}
	return batch, nil
}
