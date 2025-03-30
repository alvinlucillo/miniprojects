package controllers

import (
	"encoding/json"
	"net/http"
	"skaffoldapp/internal/api/services"
	"skaffoldapp/internal/shared/models"
)

type BatchController struct {
	batchService services.BatchService
}

func NewBatchController(batchService services.BatchService) BatchController {
	return BatchController{
		batchService: batchService,
	}
}

func (u BatchController) GetGenerateDBExportRequests(w http.ResponseWriter, r *http.Request) {
	requests, err := u.batchService.GetGenerateDBExportRequests(r.Context())
	if err != nil {
		http.Error(w, "failed fetching dbexports", http.StatusInternalServerError)
		return
	}

	var response []models.DBExportResponse
	for _, batch := range requests {
		response = append(response, models.DBExportResponse{
			ID:            batch.ID.Hex(),
			Status:        batch.Status,
			ErrorMessage:  batch.ErrorMessage,
			DateRequested: batch.DateRequested,
		})
	}

	json.NewEncoder(w).Encode(response)
}

func (u BatchController) GenerateDBExport(w http.ResponseWriter, r *http.Request) {
	dbExport, err := u.batchService.GenerateDBExport(r.Context())
	if err != nil {
		http.Error(w, "failed generating db export", http.StatusInternalServerError)
		return
	}

	response := models.DBExportResponse{
		ID:            dbExport.ID.Hex(),
		Status:        dbExport.Status,
		ErrorMessage:  dbExport.ErrorMessage,
		DateRequested: dbExport.DateRequested,
	}

	json.NewEncoder(w).Encode(response)
}
