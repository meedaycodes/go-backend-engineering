// Package model defines the core data structures used across the application.
package model

// User represents a user entity in the system.
type User struct {
	ID    string `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
}

// CreateUserRequest represents the expected payload for creating a new user.
type CreateUserRequest struct {
	Name  string `json:"name"`
	Email string `json:"email"`
}
