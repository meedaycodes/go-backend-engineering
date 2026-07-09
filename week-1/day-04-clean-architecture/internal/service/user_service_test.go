package service

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/meedaycodes/day04-clean-architecture/internal/model"
)

type MockUserRepository struct {
	mock.Mock
}

func (m *MockUserRepository) Save(user model.User) error {

	args := m.Called(user)

	return args.Error(0)
}

func (m *MockUserRepository) FindByID(ID string) (model.User, error) {
	args := m.Called(ID)
	return args.Get(0).(model.User), args.Error(1)
}

func (m *MockUserRepository) FindAll() ([]model.User, error) {
	args := m.Called()
	return args.Get(0).([]model.User), args.Error(1)
}

func TestCreateUser(t *testing.T) {

	mockRepo := new(MockUserRepository)
	mockRepo.On("Save", mock.Anything).Return(nil)

	svc := NewUserService(mockRepo)

	t.Run("Valid request", func(t *testing.T) {
		req := model.CreateUserRequest{Name: "Habeeb", Email: "habeebaramide@yahoo.com"}

		user, err := svc.CreateUser(req)

		assert.NoError(t, err)
		assert.Equal(t, user.Name, req.Name)
		assert.Equal(t, user.Email, req.Email)
		assert.NotEmpty(t, user.ID)
	})

	t.Run("empty email", func(t *testing.T) {

		req := model.CreateUserRequest{Name: "Fahma", Email: ""}

		_, err := svc.CreateUser(req)
		assert.Equal(t, errEmptyEmail, err)

	})

	t.Run("empty name", func(t *testing.T) {

		req := model.CreateUserRequest{Name: "", Email: "fahma.a.g@gmail.com"}

		_, err := svc.CreateUser(req)
		assert.Equal(t, errEmptyName, err)

	})
}

func TestGetUserByID(t *testing.T) {

	mockRepo := new(MockUserRepository)
	mockRepo.On("Save", mock.Anything).Return(nil)

	svc := NewUserService(mockRepo)

	t.Run("Valid saved user", func(t *testing.T) {

		req := model.CreateUserRequest{Name: "Habeeb", Email: "habeebaramide@yahoo.com"}

		user, err := svc.CreateUser(req)
		assert.NoError(t, err)

		mockRepo.On("FindByID", user.ID).Return(user, nil)

		getUser, err := svc.GetUserByID(user.ID)

		assert.NotEmpty(t, getUser)
		assert.Equal(t, req.Email, getUser.Email)
		assert.Equal(t, req.Name, getUser.Name)
	})

	t.Run("not Found", func(t *testing.T) {

		mockRepo.On("FindByID", "nonexistent-ID").Return(model.User{}, errors.New("User not found"))

		_, err := svc.GetUserByID("nonexistent-ID")
		assert.Error(t, err)
	})

}

func TestGetAllUsers(t *testing.T) {

	mockRepo := new(MockUserRepository)
	mockRepo.On("Save", mock.Anything).Return(nil)

	svc := NewUserService(mockRepo)

	t.Run("Get all users", func(t *testing.T) {

		var usersList []model.User

		reqs := []model.CreateUserRequest{
			{Name: "Habeeb", Email: "habeebaramide@yahoo.com"},
			{Name: "Fahma", Email: "fahma.a.gmail.com"},
		}

		for _, req := range reqs {
			user, err := svc.CreateUser(req)
			assert.NoError(t, err)

			usersList = append(usersList, user)
		}

		mockRepo.On("FindAll").Return(usersList, nil)

		getUsers, err := svc.GetAllUsers()

		assert.NoError(t, err)
		assert.NotEmpty(t, getUsers)
		assert.Equal(t, len(reqs), len(getUsers))
	})

}
