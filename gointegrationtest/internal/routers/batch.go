package routers

import (
	"gointegrationtest/internal/controllers"
	"net/http"
)

func SetupBatchRoutes(mux *http.ServeMux, batchController controllers.BatchController) {
	mux.HandleFunc("GET /generateddbexport", batchController.GetGenerateDBExportRequests)
	mux.HandleFunc("POST /generatedbexport", batchController.GenerateDBExport)
}
