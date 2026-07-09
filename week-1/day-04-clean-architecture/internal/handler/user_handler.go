// Package handler provides HTTP handlers that translate incoming requests
// into service calls and write JSON responses back to the client.
package handler

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/meedaycodes/day04-clean-architecture/internal/model"
	"github.com/meedaycodes/day04-clean-architecture/internal/service"
)

// UserHandler exposes user operations over HTTP.
type UserHandler struct {
	service *service.UserService
}

// NewUserHandler creates a new UserHandler with the given service.
func NewUserHandler(serv *service.UserService) *UserHandler {
	newServ := &UserHandler{service: serv}

	return newServ
}

// CreateUser handles POST /users. It decodes a CreateUserRequest from the body,
// delegates to the service layer, and responds with the created user as JSON.
func (u *UserHandler) CreateUser(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {

		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	var userBody model.CreateUserRequest
	dec := json.NewDecoder(r.Body)

	err := dec.Decode(&userBody)
	if err != nil {
		http.Error(w, "request body empty", http.StatusBadRequest)
		return
	}

	user, err := u.service.CreateUser(userBody)
	if err != nil {
		http.Error(w, "user not created", http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(user)

}

// GetUserByID handles GET /users/{id}. It extracts the user ID from the URL path
// and returns the matching user as JSON, or 400 if not found.
func (u *UserHandler) GetUserByID(w http.ResponseWriter, r *http.Request) {

	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	id := strings.TrimPrefix(r.URL.Path, "/users/")
	if id == "" {
		http.Error(w, "id cannot be empty", http.StatusBadRequest)
		return
	}

	user, err := u.service.GetUserByID(id)
	if err != nil {
		http.Error(w, "user not retrieved", http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(user)

}

// GetAllUsers handles GET /users. It returns all users in the system as a JSON array.
func (u *UserHandler) GetAllUsers(w http.ResponseWriter, r *http.Request) {

	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	users, err := u.service.GetAllUsers()

	if err != nil {
		http.Error(w, "users not retrieved", http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(users)

}
