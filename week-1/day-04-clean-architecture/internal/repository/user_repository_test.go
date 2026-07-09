package repository

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/meedaycodes/day04-clean-architecture/internal/model"
)

func TestInMemoryUserRepository(t *testing.T) {

	u := model.User{ID: "123ABS", Name: "Habeeb", Email: "habeebaramide@yahoo.com"}

	repo := NewInMemoryUserRepository()

	t.Run("Save and FindByID", func(t *testing.T) {

		err := repo.Save(u)
		assert.NoError(t, err)

	})

	t.Run("FindByID not found", func(t *testing.T) {

		nonUser, err := repo.FindByID("4534acd")

		assert.Equal(t, model.User{}, nonUser)
		assert.Equal(t, ErrUserNotFound, err)

	})

	t.Run("FindByID found", func(t *testing.T) {
		user, err := repo.FindByID(u.ID)

		assert.NoError(t, err)

		assert.Equal(t, u.ID, user.ID)
		assert.Equal(t, u.Email, user.Email)
		assert.Equal(t, u.Name, user.Name)

	})

	t.Run("FindAll", func(t *testing.T) {

		emptyRepo := NewInMemoryUserRepository()

		callUsers, err := emptyRepo.FindAll()

		assert.NoError(t, err)
		assert.Nil(t, callUsers)

	})
}
