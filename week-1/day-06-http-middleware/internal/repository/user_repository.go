// Package repository defines the data access layer and provides storage
// implementations for application entities.
package repository

import (
	"errors"
	"sync"

	"github.com/meedaycodes/day-06-http-middleware/internal/model"
)

// ErrUserNotFound is returned when a user with the given ID does not exist in the repository.
var (
	ErrUserNotFound = errors.New("User with the Id not found")
)

// UserRepository defines the contract for user data access operations.
// Any storage backend (in-memory, PostgreSQL, etc.) must implement this interface.
type UserRepository interface {
	// Save persists a user to the store.
	Save(user model.User) error
	// FindByID retrieves a user by their unique ID.
	FindByID(ID string) (model.User, error)
	// FindAll returns all users in the store.
	FindAll() ([]model.User, error)

	// Update overwrites an existing user in the store.
	Update(user model.User) error
	// Delete removes a user from the store by their ID.
	Delete(ID string) error
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
func (m *InMemoryUserRepository) Save(user model.User) error {
	m.mut.Lock()
	defer m.mut.Unlock()
	m.mem[user.ID] = user

	return nil
}

// FindByID returns the user with the given ID, or ErrUserNotFound if no match exists.
func (m *InMemoryUserRepository) FindByID(ID string) (model.User, error) {
	m.mut.RLock()
	defer m.mut.RUnlock()

	user, exist := m.mem[ID]

	if !exist {
		return user, ErrUserNotFound
	}

	return user, nil

}

// FindAll returns all users currently stored in memory.
func (m *InMemoryUserRepository) FindAll() (users []model.User, err error) {

	m.mut.RLock()
	defer m.mut.RUnlock()

	for _, value := range m.mem {
		users = append(users, value)
	}
	return users, nil

}

// Update replaces an existing user in the map. Returns ErrUserNotFound if the user does not exist.
func (m *InMemoryUserRepository) Update(user model.User) error {

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
func (m *InMemoryUserRepository) Delete(ID string) error {

	m.mut.Lock()
	defer m.mut.Unlock()

	_, exists := m.mem[ID]

	if !exists {
		return ErrUserNotFound
	}

	delete(m.mem, ID)

	return nil

}
