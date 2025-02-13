package routers

import (
	"gointegrationtest/internal/controllers"
	"net/http"
)

func SetupBatchRoutes(mux *http.ServeMux, batchController controllers.BatchController) {
	mux.HandleFunc("GET /batches", batchController.GetBatches)
	mux.HandleFunc("POST /batches", batchController.CreateBatch)
}
