package routers

import (
	"net/http"
	"skaffoldapp/internal/api/controllers"
)

func SetupBatchRoutes(mux *http.ServeMux, batchController controllers.BatchController) {
	mux.HandleFunc("GET /generateddbexport", batchController.GetGenerateDBExportRequests)
	mux.HandleFunc("POST /generatedbexport", batchController.GenerateDBExport)
}
