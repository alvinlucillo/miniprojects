package controllers

import (
	"encoding/json"
	"gointegrationtest/internal/models"
	"gointegrationtest/internal/services"
	"net/http"
)

type BatchController struct {
	batchService services.BatchService
}

func NewBatchController(batchService services.BatchService) BatchController {
	return BatchController{
		batchService: batchService,
	}
}

func (u BatchController) GetBatches(w http.ResponseWriter, r *http.Request) {
	batches, err := u.batchService.GetBatches(r.Context())
	if err != nil {
		http.Error(w, "failed to fetch batches", http.StatusInternalServerError)
		return
	}

	var response []models.BatchResponse
	for _, batch := range batches {
		response = append(response, models.BatchResponse{
			ID:           batch.ID.Hex(),
			Status:       batch.Status,
			ErrorMessage: batch.ErrorMessage,
		})
	}

	json.NewEncoder(w).Encode(response)
}

func (u BatchController) CreateBatch(w http.ResponseWriter, r *http.Request) {
	id, err := u.batchService.CreateBatch(r.Context())
	if err != nil {
		http.Error(w, "failed to create batch", http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(models.BatchResponse{
		ID: id,
	})
}
