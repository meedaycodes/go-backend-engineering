// Package handler provides HTTP handlers that translate incoming requests
// into service calls and write JSON responses back to the client.
package handler

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"

	"github.com/meedaycodes/day-06-http-middleware/internal/model"
	"github.com/meedaycodes/day-06-http-middleware/internal/service"
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

	id := chi.URLParam(r, "id")
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

	users, err := u.service.GetAllUsers()

	if err != nil {
		http.Error(w, "users not retrieved", http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(users)

}

// UpdateUser handles PUT /users/{id}. It extracts the user ID from the URL path,
// decodes an UpdateUserRequest from the body, and responds with the updated user as JSON.
func (u *UserHandler) UpdateUser(w http.ResponseWriter, r *http.Request) {

	id := chi.URLParam(r, "id")
	dec := json.NewDecoder(r.Body)

	var userBody model.UpdateUserRequest

	err := dec.Decode(&userBody)
	if err != nil {
		http.Error(w, "request body empty", http.StatusBadRequest)
		return
	}

	newUser, err := u.service.UpdateUser(id, userBody)

	if err != nil {
		http.Error(w, "user not updated", http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(newUser)

}

// DeleteUser handles DELETE /users/{id}. It removes the user and responds with 204 No Content.
func (u *UserHandler) DeleteUser(w http.ResponseWriter, r *http.Request) {

	id := chi.URLParam(r, "id")

	err := u.service.DeleteUser(id)
	if err != nil {
		http.Error(w, "delete user failed", http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusNoContent)

}
