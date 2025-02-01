package routers

import (
	"gobddintgtest/internal/controllers"
	"net/http"
)

func SetupUserRoutes(mux *http.ServeMux, userController controllers.UsersController) {
	mux.HandleFunc("GET /users", userController.GetUsers)
}
