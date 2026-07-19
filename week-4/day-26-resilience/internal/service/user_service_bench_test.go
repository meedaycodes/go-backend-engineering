package service

import (
	"context"
	"fmt"
	"testing"

	"github.com/meedaycodes/day26-resilience/internal/model"
	"github.com/meedaycodes/day26-resilience/internal/repository"
)

func BenchmarkCreateUser(b *testing.B) {

	repo := repository.NewInMemoryUserRepository()
	svc := NewUserService(repo, nil)
	b.ResetTimer()

	for i := 0; i < b.N; i++ {

		svc.CreateUser(context.Background(), model.CreateUserRequest{Name: "Habeeb", Email: fmt.Sprintf("user%d@test.com", i), Password: "password"}) //nolint:errcheck

	}
}
