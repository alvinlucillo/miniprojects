package repos

import (
	"context"
	"fmt"
	"skaffoldapp/internal/shared/database"
	"skaffoldapp/internal/shared/models"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type BatchRepo struct {
	exportedDBColl *mongo.Collection
}

func NewBatchRepo(db *mongo.Database) BatchRepo {
	return BatchRepo{
		exportedDBColl: db.Collection(database.EXPORTED_DB_COLLECTION),
	}
}

func (b BatchRepo) GetDBExports(ctx context.Context) ([]models.DBExport, error) {
	var dbExports []models.DBExport

	if err := database.FindAll(b.exportedDBColl, &dbExports); err != nil {
		return nil, fmt.Errorf("failed to get batches: %w", err)
	}
	return dbExports, nil
}

func (b BatchRepo) InsertDBExports(ctx context.Context, dbExport models.DBExport) (primitive.ObjectID, error) {
	result, err := b.exportedDBColl.InsertOne(ctx, dbExport)
	if err != nil {
		return primitive.NilObjectID, fmt.Errorf("failed to insert dbExport: %w", err)
	}

	id := result.InsertedID.(primitive.ObjectID)
	return id, nil
}

func (b BatchRepo) UpdateDBExport(ctx context.Context, dbExport models.DBExport) error {
	_, err := b.exportedDBColl.UpdateOne(ctx, bson.M{"_id": dbExport.ID}, bson.M{"$set": dbExport})
	if err != nil {
		return fmt.Errorf("failed to update dbExport: %w", err)
	}
	return nil
}
