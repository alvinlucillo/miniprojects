package routers

import (
	"gointegrationtest/internal/controllers"
	"net/http"
)

func SetupUserRoutes(mux *http.ServeMux, userController controllers.UsersController) {
	mux.HandleFunc("GET /users", userController.GetUsers)
	mux.HandleFunc("POST /users", userController.CreateUser)
}
