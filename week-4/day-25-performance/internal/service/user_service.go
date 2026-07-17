// Package service contains the business logic for the application.
// It validates input, enforces rules, and delegates persistence to the
// repository layer. The user service also integrates a Redis cache using the
// cache-aside pattern: reads check the cache first, writes invalidate it.
package service

import (
	"context"
	"encoding/json"
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/meedaycodes/day25-performance/internal/cache"
	"github.com/meedaycodes/day25-performance/internal/model"
	"github.com/meedaycodes/day25-performance/internal/repository"
)

var (
	errEmptyName  = errors.New("name cannot be empty")
	errEmptyEmail = errors.New("email cannot be empty")
)

// UserService provides business operations for managing users.
// It depends on the UserRepository interface (not a concrete implementation)
// and a RedisCache for caching individual user lookups.
type UserService struct {
	repo  repository.UserRepository
	cache *cache.RedisCache
}

// NewUserService creates a new UserService with the given repository and Redis
// cache. Both are required — the cache is used in the read path and for
// invalidation on writes.
func NewUserService(repo repository.UserRepository, cache *cache.RedisCache) *UserService {
	newRepo := &UserService{repo: repo, cache: cache}

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

// GetUserByID implements the cache-aside pattern. It checks Redis first using
// the key "user:<ID>". On a cache hit, the JSON value is unmarshalled and
// returned immediately — no DB call. On a miss (any error from cache.Get),
// it falls through to the repository, marshals the result to JSON, stores it
// with a 15-minute TTL, then returns the user. The TTL caps how long stale
// data can linger if invalidation is missed.
func (s *UserService) GetUserByID(ctx context.Context, ID string) (user model.User, err error) {

	cacheKey := "user:" + ID

	val, err := s.cache.Get(ctx, cacheKey)
	if err != nil {

		user, err = s.repo.FindByID(ctx, ID)
		if err != nil {
			return user, err
		}

		marshUser, err := json.Marshal(user)
		if err != nil {
			return user, err
		}

		_ = s.cache.Set(ctx, cacheKey, string(marshUser), 15*time.Minute)
		return user, nil
	}

	if err = json.Unmarshal([]byte(val), &user); err != nil {
		return user, err
	}

	return user, nil
}

// GetAllUsers returns all users in the system with pagination.
func (s *UserService) GetAllUsers(ctx context.Context, limit, offset int) (allUsers []model.User, err error) {

	allUsers, err = s.repo.FindAll(ctx, limit, offset)

	return allUsers, err
}

// UpdateUser validates the request, looks up the existing user, applies
// changes, and persists the update. After a successful repo write, the cache
// entry for this user is invalidated so subsequent reads fetch fresh data.
// Returns errEmptyName or errEmptyEmail if required fields are missing, or
// ErrUserNotFound if the user does not exist.
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
	_ = s.cache.Delete(ctx, "user:"+ID)
	return findUser, nil
}

// DeleteUser removes a user by their ID. After a successful repo delete, the
// cache entry is invalidated so the key is not served after the user is gone.
// Returns ErrUserNotFound if the user does not exist.
func (s *UserService) DeleteUser(ctx context.Context, ID string) error {

	err := s.repo.Delete(ctx, ID)
	if err != nil {
		return err
	}
	_ = s.cache.Delete(ctx, "user:"+ID)
	return nil
}
