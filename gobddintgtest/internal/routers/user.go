package routers

import (
	"gobddintgtest/internal/controllers"
	"net/http"
)

func SetupUserRoutes(userController controllers.UsersController) *http.ServeMux {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /users", userController.GetUsers)

	return mux
}
