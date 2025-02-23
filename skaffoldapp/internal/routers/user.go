package routers

import (
	"net/http"
	"skaffoldapp/internal/controllers"
)

func SetupUserRoutes(mux *http.ServeMux, userController controllers.UsersController) {
	mux.HandleFunc("GET /users", userController.GetUsers)
	mux.HandleFunc("POST /users", userController.CreateUser)
}
