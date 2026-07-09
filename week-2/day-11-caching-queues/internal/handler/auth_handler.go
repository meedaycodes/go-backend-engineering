// Package handler provides HTTP handlers that translate incoming requests
// into service calls and write JSON responses back to the client.
// auth_handler.go handles public authentication endpoints (signup, login)
// that do not require JWT authentication.
package handler

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/meedaycodes/day11-caching-queues/internal/model"
	"github.com/meedaycodes/day11-caching-queues/internal/service"
	"github.com/meedaycodes/day11-caching-queues/internal/validator"
)

// AuthHandler exposes authentication operations over HTTP.
type AuthHandler struct {
	service *service.AuthService
}

// NewAuthHandler creates a new AuthHandler with the given auth service.
func NewAuthHandler(serv *service.AuthService) *AuthHandler {

	newServe := &AuthHandler{service: serv}

	return newServe
}

// Signup handles POST /auth/signup. Decodes a CreateUserRequest, delegates
// to the auth service, and returns tokens on success. Returns 409 Conflict
// if the email is already registered, 400 for other failures.
func (a *AuthHandler) Signup(w http.ResponseWriter, r *http.Request) {

	var userBody model.CreateUserRequest
	var vldError *validator.ValidationError
	dec := json.NewDecoder(r.Body)

	ctx := r.Context()

	err := dec.Decode(&userBody)
	if err != nil {
		http.Error(w, "request body empty", http.StatusBadRequest)
		return
	}

	err = validator.ValidateCreateUser(userBody)
	if errors.As(err, &vldError) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(vldError.Fields)
		return
	}

	res, err := a.service.Signup(ctx, userBody)
	if errors.Is(err, service.ErrEmailAlreadyExists) {
		http.Error(w, "Email already exists", http.StatusConflict)
		return
	}
	if err != nil {
		http.Error(w, "signup failed", http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(res)
}

// Login handles POST /auth/login. Decodes a LoginRequest, verifies credentials
// through the auth service, and returns tokens on success. Returns 401 Unauthorized
// for invalid credentials — never reveals whether email or password was wrong.
func (a *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {

	var userBody model.LoginRequest
	var vldError *validator.ValidationError
	dec := json.NewDecoder(r.Body)

	ctx := r.Context()

	err := dec.Decode(&userBody)
	if err != nil {
		http.Error(w, "request body empty", http.StatusBadRequest)
		return
	}

	err = validator.ValidateLoginRequest(userBody)
	if errors.As(err, &vldError) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(vldError.Fields)
		return
	}

	res, err := a.service.Login(ctx, userBody)

	if errors.Is(err, service.ErrInvalidCredentials) {
		http.Error(w, "Invalid credentials provided", http.StatusUnauthorized)
		return
	}

	if err != nil {
		http.Error(w, "login unsuccessful", http.StatusBadRequest)
		return

	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(res)
}
