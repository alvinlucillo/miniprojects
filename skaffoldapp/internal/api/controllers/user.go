package controllers

import (
	"encoding/json"
	"net/http"
	"skaffoldapp/internal/api/services"
	"skaffoldapp/internal/shared/models"
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
	users, err := u.userService.GetUsers(r.Context())
	if err != nil {
		http.Error(w, "failed fetching users", http.StatusInternalServerError)
		return
	}

	response := []models.UserResponse{}
	for _, user := range users {
		response = append(response, models.UserResponse{
			ID:   user.ID.Hex(),
			Name: user.Name,
		})
	}

	json.NewEncoder(w).Encode(response)
}

func (u UsersController) CreateUser(w http.ResponseWriter, r *http.Request) {
	var req models.UserRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	id, err := u.userService.CreateUser(r.Context(), req)
	if err != nil {
		http.Error(w, "Failed to create user", http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(models.UserResponse{
		ID:   id,
		Name: req.Name,
	})
}
