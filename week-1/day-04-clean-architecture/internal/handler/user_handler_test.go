package handler

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/meedaycodes/day04-clean-architecture/internal/model"
	"github.com/meedaycodes/day04-clean-architecture/internal/service"
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
	svc := service.NewUserService(mockRepo)
	handler := NewUserHandler(svc)

	mockRepo.On("Save", mock.Anything).Return(nil)

	t.Run("Save correct request", func(t *testing.T) {

		req := httptest.NewRequest(http.MethodPost, "/users", strings.NewReader(`{"name": "habeeb", "email":"habeebaramide@uyahoo.com"}`))
		rec := httptest.NewRecorder()

		handler.CreateUser(rec, req)

		assert.Equal(t, rec.Code, http.StatusCreated)
	})

	t.Run("Save Empty Body Request", func(t *testing.T) {

		req := httptest.NewRequest(http.MethodPost, "/users", strings.NewReader(``))
		rec := httptest.NewRecorder()

		handler.CreateUser(rec, req)

		assert.Equal(t, rec.Code, http.StatusBadRequest)
	})

	t.Run("Wrong Method", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/users", strings.NewReader(`{"name": "habeeb", "email":"habeebaramide@uyahoo.com"}`))
		rec := httptest.NewRecorder()

		handler.CreateUser(rec, req)

		assert.Equal(t, rec.Code, http.StatusMethodNotAllowed)

	})

}

func TestGetUserByID(t *testing.T) {

	mockRepo := new(MockUserRepository)
	svc := service.NewUserService(mockRepo)
	handler := NewUserHandler(svc)

	mockRepo.On("Save", mock.Anything).Return(nil)

	t.Run("Retrieve Saved user by ID", func(t *testing.T) {

		expectedUser := model.User{ID: "abc-1234", Name: "habeeb", Email: "habeebaramide@yahoo.com"}
		mockRepo.On("FindByID", "abc-1234").Return(expectedUser, nil)

		req := httptest.NewRequest(http.MethodGet, "/users/abc-1234", nil)
		rec := httptest.NewRecorder()

		handler.GetUserByID(rec, req)

		assert.Equal(t, http.StatusOK, rec.Code)

	})

	t.Run("Retrieve with empty ID", func(t *testing.T) {

		mockRepo.On("FindByID", "").Return(model.User{}, errors.New("ID provided is empty"))

		req := httptest.NewRequest(http.MethodGet, "/users/", nil)
		rec := httptest.NewRecorder()

		handler.GetUserByID(rec, req)

		assert.Equal(t, http.StatusBadRequest, rec.Code)

	})

	t.Run("Wrong Method", func(t *testing.T) {

		expectedUser := model.User{ID: "abc-1234", Name: "habeeb", Email: "habeebaramide@yahoo.com"}
		mockRepo.On("FindByID", "abc-1234").Return(expectedUser, nil)

		req := httptest.NewRequest(http.MethodPost, "/users/abc-1234", nil)
		rec := httptest.NewRecorder()

		handler.GetUserByID(rec, req)

		assert.Equal(t, rec.Code, http.StatusMethodNotAllowed)

	})
}

func TestGetAllUsers(t *testing.T) {

	mockRepo := new(MockUserRepository)
	svc := service.NewUserService(mockRepo)
	handler := NewUserHandler(svc)

	mockRepo.On("Save", mock.Anything).Return(nil)

	t.Run("Get all Users Saved", func(t *testing.T) {

		expectedUsers := []model.User{
			{ID: "aq-123", Name: "Habeeb", Email: "habeebaramide@yahoo.com"},
			{ID: "aq-124", Name: "Fahma", Email: "habeebamide@yahoo.com"},
		}

		mockRepo.On("FindAll").Return(expectedUsers, nil)

		req := httptest.NewRequest(http.MethodGet, "/users", nil)
		rec := httptest.NewRecorder()

		handler.GetAllUsers(rec, req)
		assert.Equal(t, http.StatusOK, rec.Code)
	})

}
