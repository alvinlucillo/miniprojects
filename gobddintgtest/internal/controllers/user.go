package controllers

import (
	"encoding/json"
	"gobddintgtest/internal/services"
	"net/http"
)

type UsersController struct {
	userService services.UserService
}

func NewUsersController(userService services.UserService) UsersController {
	return UsersController{
		userService: userService,
	}
}

func (u UsersController) GetUsers(w http.ResponseWriter, r *http.Request) {
	// users, err := u.userService.GetUsers(r.Context())
	// if err != nil {
	// 	http.Error(w, "Failed to fetch users", http.StatusInternalServerError)
	// 	return
	// }
	// json.NewEncoder(w).Encode(users)

	json.NewEncoder(w).Encode("henlo world")
}

/**
func (c *UserController) GetUsers(w http.ResponseWriter, r *http.Request) {
	users, err := c.service.GetUsers(r.Context())
	if err != nil {
		http.Error(w, "Failed to fetch users", http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(users)
}
	**/
