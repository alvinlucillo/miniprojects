package controllers

import (
	"encoding/json"
	"gointegrationtest/internal/services"
	"net/http"
)

type UserResponse struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

type UsersController struct {
	userService services.UserService
}

func NewUsersController(userService services.UserService) UsersController {
	return UsersController{
		userService: userService,
	}
}

func (u UsersController) GetUsers(w http.ResponseWriter, r *http.Request) {
	users, err := u.userService.GetUsers(r.Context())
	if err != nil {
		http.Error(w, "Failed to fetch users", http.StatusInternalServerError)
		return
	}

	var response []UserResponse
	for _, user := range users {
		response = append(response, UserResponse{
			ID:   user.ID.Hex(),
			Name: user.Name,
		})
	}

	json.NewEncoder(w).Encode(response)
}
