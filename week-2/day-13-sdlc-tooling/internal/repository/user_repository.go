// Package repository defines the data access layer and provides storage
// implementations for application entities.
package repository

import (
	"context"
	"errors"
	"sync"

	"github.com/meedaycodes/day13-sdlc-tooling/internal/model"
)

// ErrUserNotFound is returned when a user with the given ID does not exist in the repository.
var (
	ErrUserNotFound = errors.New("user with the Id not found")
)

// UserRepository defines the contract for user data access operations.
// Any storage backend (in-memory, PostgreSQL, etc.) must implement this interface.
type UserRepository interface {
	// Save persists a user to the store.
	Save(ctx context.Context, user model.User) error
	// FindByID retrieves a user by their unique ID.
	FindByID(ctx context.Context, ID string) (model.User, error)
	// FindAll returns all users in the store.
	FindAll(ctx context.Context, limit, offset int) ([]model.User, error)

	// Update overwrites an existing user in the store.
	Update(ctx context.Context, user model.User) error
	// Delete removes a user from the store by their ID.
	Delete(ctx context.Context, ID string) error

	FindByEmail(ctx context.Context, email string) (model.User, error)
}

// InMemoryUserRepository is a thread-safe, map-based implementation of UserRepository.
// It is intended for development and testing, not production use.
type InMemoryUserRepository struct {
	mem map[string]model.User
	mut sync.RWMutex
}

// NewInMemoryUserRepository creates and returns a new InMemoryUserRepository.
func NewInMemoryUserRepository() *InMemoryUserRepository {
	var InMemoryUserRepository InMemoryUserRepository
	InMemoryUserRepository.mem = make(map[string]model.User, 0)

	return &InMemoryUserRepository
}

// Save stores a user in the in-memory map, keyed by their ID.
func (m *InMemoryUserRepository) Save(ctx context.Context, user model.User) error {
	m.mut.Lock()
	defer m.mut.Unlock()
	m.mem[user.ID] = user

	return nil
}

// FindByID returns the user with the given ID, or ErrUserNotFound if no match exists.
func (m *InMemoryUserRepository) FindByID(ctx context.Context, ID string) (model.User, error) {
	m.mut.RLock()
	defer m.mut.RUnlock()

	user, exist := m.mem[ID]

	if !exist {
		return user, ErrUserNotFound
	}

	return user, nil

}

// FindAll returns all users currently stored in memory.
func (m *InMemoryUserRepository) FindAll(ctx context.Context, limit, offset int) (users []model.User, err error) {

	var slicedUsers []model.User

	m.mut.RLock()
	defer m.mut.RUnlock()

	for _, value := range m.mem {
		users = append(users, value)
	}

	if offset >= len(users) {
		return slicedUsers, nil

	} else if offset+limit >= len(users) {
		slicedUsers = users[offset:]
		return slicedUsers, nil
	} else {
		slicedUsers = users[offset : offset+limit]
		return slicedUsers, nil
	}

}

// Update replaces an existing user in the map. Returns ErrUserNotFound if the user does not exist.
func (m *InMemoryUserRepository) Update(ctx context.Context, user model.User) error {

	m.mut.Lock()
	defer m.mut.Unlock()

	_, exist := m.mem[user.ID]

	if !exist {
		return ErrUserNotFound
	}

	m.mem[user.ID] = user

	return nil

}

// Delete removes a user from the map by their ID. Returns ErrUserNotFound if the user does not exist.
func (m *InMemoryUserRepository) Delete(ctx context.Context, ID string) error {

	m.mut.Lock()
	defer m.mut.Unlock()

	_, exists := m.mem[ID]

	if !exists {
		return ErrUserNotFound
	}

	delete(m.mem, ID)

	return nil

}

// Finds a user with an existing email passed into the payload
func (m *InMemoryUserRepository) FindByEmail(ctx context.Context, email string) (user model.User, err error) {

	m.mut.RLock()
	defer m.mut.RUnlock()

	for _, user := range m.mem {

		if user.Email == email {
			return user, nil
		}

	}

	return user, ErrUserNotFound
}
