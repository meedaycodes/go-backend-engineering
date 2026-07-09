// Package service contains the business logic for the application.
// It validates input, enforces rules, and delegates persistence to the repository layer.
package service

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/meedaycodes/day10-api-design/internal/model"
	"github.com/meedaycodes/day10-api-design/internal/repository"
)

var (
	errEmptyName  = errors.New("Name cannot be empty")
	errEmptyEmail = errors.New("Email cannot be empty")
)

// UserService provides business operations for managing users.
// It depends on the UserRepository interface, not a concrete implementation.
type UserService struct {
	repo repository.UserRepository
}

// NewUserService creates a new UserService with the given repository.
func NewUserService(repo repository.UserRepository) *UserService {
	newRepo := &UserService{repo: repo}

	return newRepo
}

// CreateUser validates the request, generates a UUID, and persists the new user.
// It returns errEmptyName or errEmptyEmail if required fields are missing.
func (s *UserService) CreateUser(ctx context.Context, r model.CreateUserRequest) (user model.User, err error) {

	if r.Name == "" {
		return user, errEmptyName
	}
	if r.Email == "" {
		return user, errEmptyEmail
	}

	id := uuid.New().String()
	user = model.User{ID: id, Name: r.Name, Email: r.Email}

	saveErr := s.repo.Save(ctx, user)
	if saveErr != nil {
		return user, saveErr
	}

	return user, nil

}

// GetUserByID retrieves a user by their ID. Returns repository.ErrUserNotFound if the user does not exist.
func (s *UserService) GetUserByID(ctx context.Context, ID string) (user model.User, err error) {
	user, err = s.repo.FindByID(ctx, ID)
	return user, err
}

// GetAllUsers returns all users in the system.
func (s *UserService) GetAllUsers(ctx context.Context, limit, offset int) (allUsers []model.User, err error) {

	allUsers, err = s.repo.FindAll(ctx, limit, offset)

	return allUsers, err
}

// UpdateUser validates the request, looks up the existing user, applies changes, and persists the update.
// Returns errEmptyName or errEmptyEmail if required fields are missing, or ErrUserNotFound if the user does not exist.
func (s *UserService) UpdateUser(ctx context.Context, ID string, r model.UpdateUserRequest) (user model.User, err error) {

	findUser, err := s.repo.FindByID(ctx, ID)

	if err != nil {
		return findUser, err
	}

	if r.Name == "" {
		return user, errEmptyName
	}

	if r.Email == "" {
		return user, errEmptyEmail
	}

	findUser.Name = r.Name
	findUser.Email = r.Email

	err = s.repo.Update(ctx, findUser)

	if err != nil {
		return findUser, err
	}

	return findUser, nil
}

// DeleteUser removes a user by their ID. Returns ErrUserNotFound if the user does not exist.
func (s *UserService) DeleteUser(ctx context.Context, ID string) error {

	err := s.repo.Delete(ctx, ID)
	if err != nil {
		return err
	}

	return nil
}
