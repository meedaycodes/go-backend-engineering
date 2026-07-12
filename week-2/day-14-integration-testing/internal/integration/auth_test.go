package integration_test

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"

	"github.com/golang-jwt/jwt/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
	"github.com/testcontainers/testcontainers-go/modules/redis"
	"github.com/testcontainers/testcontainers-go/wait"

	"github.com/go-chi/chi/v5"
	"github.com/meedaycodes/day14-integration-testing/internal/cache"
	"github.com/meedaycodes/day14-integration-testing/internal/handler"
	"github.com/meedaycodes/day14-integration-testing/internal/middleware"
	"github.com/meedaycodes/day14-integration-testing/internal/model"
	"github.com/meedaycodes/day14-integration-testing/internal/repository"
	"github.com/meedaycodes/day14-integration-testing/internal/service"
	"github.com/meedaycodes/day14-integration-testing/internal/worker"
)

var (
	testRouter http.Handler
	testDB     *pgxpool.Pool
)

const JWTSecret = "test-secret"

func TestMain(m *testing.M) {

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	pgContainer, err := postgres.Run(ctx, "postgres:16-alpine",
		postgres.WithDatabase("testdb"),
		postgres.WithUsername("test"),
		postgres.WithPassword("test"),
		testcontainers.WithWaitStrategy(
			wait.ForListeningPort("5432/tcp"),
		),
	)
	if err != nil {
		log.Fatal(err)
	}

	connStr, err := pgContainer.ConnectionString(ctx, "sslmode=disable")
	if err != nil {
		log.Fatal(err)
	}

	testDB, err = pgxpool.New(ctx, connStr)
	if err != nil {
		log.Fatal(err)
	}

	sqlFile, err := os.ReadFile("../../migrations/000001_create_users_table.up.sql")
	if err != nil {
		log.Fatal(err)
	}

	_, err = testDB.Exec(ctx, string(sqlFile))

	if err != nil {
		log.Fatal(err)
	}

	redisContainer, err := redis.Run(ctx, "redis:7-alpine")
	if err != nil {
		log.Fatal(err)
	}

	redisAddr, err := redisContainer.Endpoint(ctx, "")
	if err != nil {
		log.Fatal(err)
	}

	testEmailWorker := worker.NewEmailWorker(100)
	go testEmailWorker.Start(ctx)

	repo := repository.NewPostgresUserRepository(testDB)
	testCache := cache.NewRedisCache(redisAddr)
	testSvc := service.NewUserService(repo, testCache)
	testAuthSvc := service.NewAuthService(repo, JWTSecret, testEmailWorker)
	testUserHandler := handler.NewUserHandler(testSvc)
	testAuthHandler := handler.NewAuthHandler(testAuthSvc)

	r := chi.NewRouter()

	r.Use(middleware.RateLimit())
	r.Use(middleware.Recover)
	r.Use(middleware.Logging)

	r.Post("/auth/signup", testAuthHandler.Signup)
	r.Post("/auth/login", testAuthHandler.Login)

	r.Route("/users", func(r chi.Router) {

		r.Use(middleware.Auth(JWTSecret))
		r.Get("/", testUserHandler.GetAllUsers)
		r.Get("/{id}", testUserHandler.GetUserByID)
		r.Put("/{id}", testUserHandler.UpdateUser)
		r.Delete("/{id}", testUserHandler.DeleteUser)
	})

	testRouter = r

	os.Exit(m.Run())
}

func TestSignup(t *testing.T) {

	tests := []struct {
		name       string
		body       string
		wantStatus int
	}{
		{"happy path", `{"name":"John","email":"john@test.com","password":"password123"}`, 201},
		{"duplicate email", `{"name":"John","email":"john@test.com","password":"password123"}`, 409},
		{"invalid body", `{}`, 400},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodPost, "/auth/signup", strings.NewReader(tt.body))
			req.Header.Set("Content-Type", "application/json")
			rec := httptest.NewRecorder()
			testRouter.ServeHTTP(rec, req)
			if rec.Code != tt.wantStatus {
				t.Errorf("got %d, want %d", rec.Code, tt.wantStatus)
			}
		})
	}
	t.Cleanup(func() {
		_, err := testDB.Exec(context.Background(), "DELETE FROM users")
		if err != nil {
			t.Fatal(err)
		}
	})
}

func TestLogin(t *testing.T) {

	// setup: create a user for login tests to use
	body := strings.NewReader(`{"name":"John","email":"john@test.com","password":"password123"}`)
	req := httptest.NewRequest(http.MethodPost, "/auth/signup", body)
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	testRouter.ServeHTTP(rec, req)

	tests := []struct {
		name       string
		body       string
		wantStatus int
	}{
		{"happy path", `{"email":"john@test.com","password":"password123"}`, 200},
		{"wrong password", `{"email":"john@test.com","password":"passwor123"}`, 401},
		{"user not found", `{"email":"john@tet.com","password":"passwor123"}`, 401},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodPost, "/auth/login", strings.NewReader(tt.body))
			req.Header.Set("Content-Type", "application/json")
			rec := httptest.NewRecorder()
			testRouter.ServeHTTP(rec, req)
			if rec.Code != tt.wantStatus {
				t.Errorf("got %d, want %d", rec.Code, tt.wantStatus)
			}
		})
	}
	t.Cleanup(func() {
		_, err := testDB.Exec(context.Background(), "DELETE FROM users")
		if err != nil {
			t.Fatal(err)
		}
	})
}

func loginUser(t *testing.T, email, password string) string {

	var res model.AuthResponse

	body := fmt.Sprintf(`{"email":"%s","password":"%s"}`, email, password)

	req := httptest.NewRequest(http.MethodPost, "/auth/login", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")

	rec := httptest.NewRecorder()
	testRouter.ServeHTTP(rec, req)
	if rec.Code != http.StatusOK {
		t.Fatalf("loginUser: expected 200, got %d", rec.Code)
	}

	decRes := json.NewDecoder(rec.Body)

	if err := decRes.Decode(&res); err != nil {
		t.Fatalf("loginUser: failed to decode response: %v", err)
	}
	return res.AccessToken

}

func TestGetUsers(t *testing.T) {

	body := strings.NewReader(`{"name":"John","email":"john@test.com","password":"password123"}`)
	req := httptest.NewRequest(http.MethodPost, "/auth/signup", body)
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	testRouter.ServeHTTP(rec, req)

	accessToken := loginUser(t, "john@test.com", "password123")

	tests := []struct {
		name       string
		token      string
		wantStatus int
	}{
		{"no auth", "", 401},
		{"with valid token", accessToken, 200},
	}

	for _, tt := range tests {

		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodGet, "/users", nil)
			req.Header.Set("Content-Type", "application/json")
			if tt.token != "" {
				req.Header.Set("Authorization", "Bearer "+tt.token)
			}
			rec := httptest.NewRecorder()
			testRouter.ServeHTTP(rec, req)
			if rec.Code != tt.wantStatus {
				t.Errorf("got %d, want %d", rec.Code, tt.wantStatus)
			}
		})
	}
	t.Cleanup(func() {
		_, err := testDB.Exec(context.Background(), "DELETE FROM users")
		if err != nil {
			t.Fatal(err)
		}
	})

}

func getUserIDFromToken(t *testing.T, tokenString string) string {

	token, _, err := new(jwt.Parser).ParseUnverified(tokenString, jwt.MapClaims{})
	if err != nil {
		log.Fatal(err)
	}
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		t.Fatal("invalid claims")
	}
	userID, _ := claims.GetSubject()
	return userID
}

func TestGetUserByID(t *testing.T) {

	body := strings.NewReader(`{"name":"John","email":"john@test.com","password":"password123"}`)
	req := httptest.NewRequest(http.MethodPost, "/auth/signup", body)
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	testRouter.ServeHTTP(rec, req)

	accessToken := loginUser(t, "john@test.com", "password123")
	userID := getUserIDFromToken(t, accessToken)

	tests := []struct {
		name       string
		ID         string
		wantStatus int
	}{
		{"valid ID", userID, 200},
		{"non-existtent ID", "GTR432", 400},
	}

	for _, tt := range tests {

		t.Run(tt.name, func(t *testing.T) {

			req := httptest.NewRequest(http.MethodGet, "/users/"+tt.ID, nil)
			req.Header.Set("Content-Type", "application/json")
			if accessToken != "" {
				req.Header.Set("Authorization", "Bearer "+accessToken)
			}
			rec := httptest.NewRecorder()
			testRouter.ServeHTTP(rec, req)
			if rec.Code != tt.wantStatus {
				t.Errorf("got %d, want %d", rec.Code, tt.wantStatus)
			}
		})
	}
	t.Cleanup(func() {
		_, err := testDB.Exec(context.Background(), "DELETE FROM users")
		if err != nil {
			t.Fatal(err)
		}
	})

}

func TestUpdateUser(t *testing.T) {

	body := strings.NewReader(`{"name":"John","email":"john@test.com","password":"password123"}`)
	req := httptest.NewRequest(http.MethodPost, "/auth/signup", body)
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	testRouter.ServeHTTP(rec, req)

	accessToken := loginUser(t, "john@test.com", "password123")
	userID := getUserIDFromToken(t, accessToken)

	tests := []struct {
		name       string
		ID         string
		body       string
		token      string
		wantStatus int
	}{
		{"no auth", userID, "", "", 401},
		{"happy path", userID, `{"name":"Jane","email":"jane@test.com"}`, accessToken, 200},
	}

	for _, tt := range tests {

		t.Run(tt.name, func(t *testing.T) {

			req := httptest.NewRequest(http.MethodPut, "/users/"+tt.ID, strings.NewReader(tt.body))
			req.Header.Set("Content-Type", "application/json")
			if tt.token != "" {
				req.Header.Set("Authorization", "Bearer "+tt.token)
			}
			rec := httptest.NewRecorder()
			testRouter.ServeHTTP(rec, req)
			if rec.Code != tt.wantStatus {
				t.Errorf("got %d, want %d", rec.Code, tt.wantStatus)
			}
		})
	}
	t.Cleanup(func() {
		_, err := testDB.Exec(context.Background(), "DELETE FROM users")
		if err != nil {
			t.Fatal(err)
		}
	})

}

func TestDeleteUser(t *testing.T) {

	body := strings.NewReader(`{"name":"John","email":"john@test.com","password":"password123"}`)
	req := httptest.NewRequest(http.MethodPost, "/auth/signup", body)
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	testRouter.ServeHTTP(rec, req)

	accessToken := loginUser(t, "john@test.com", "password123")
	userID := getUserIDFromToken(t, accessToken)

	tests := []struct {
		name       string
		ID         string
		token      string
		wantStatus int
	}{
		{"no auth", userID, "", 401},
		{"happy path", userID, accessToken, 204},
	}

	for _, tt := range tests {

		t.Run(tt.name, func(t *testing.T) {

			req := httptest.NewRequest(http.MethodDelete, "/users/"+tt.ID, nil)
			req.Header.Set("Content-Type", "application/json")
			if tt.token != "" {
				req.Header.Set("Authorization", "Bearer "+tt.token)
			}
			rec := httptest.NewRecorder()
			testRouter.ServeHTTP(rec, req)
			if rec.Code != tt.wantStatus {
				t.Errorf("got %d, want %d", rec.Code, tt.wantStatus)
			}
		})
	}
	t.Cleanup(func() {
		_, err := testDB.Exec(context.Background(), "DELETE FROM users")
		if err != nil {
			t.Fatal(err)
		}
	})

}
