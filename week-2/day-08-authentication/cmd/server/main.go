// Package main is the entry point that wires all layers together and starts
// the HTTP server. Auth routes (/auth/signup, /auth/login) are public.
// User routes (/users) are protected behind JWT middleware.
// Dependencies flow: pool → repo → services → handlers → router → server.
package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/meedaycodes/day08-authentication/internal/handler"
	"github.com/meedaycodes/day08-authentication/internal/middleware"
	"github.com/meedaycodes/day08-authentication/internal/repository"
	"github.com/meedaycodes/day08-authentication/internal/service"
)

func main() {

	dbURL := "postgres://habeebaramideshomuyiwa@localhost:5432/day08_users?sslmode=disable"

	pool, err := pgxpool.New(context.Background(), dbURL)

	if err != nil {
		log.Fatal(err)
	}
	log.Println("Connected to database")

	defer pool.Close()

	jwtSecret := "my-secret-string"

	repo := repository.NewPostgresUserRepository(pool)
	userSvc := service.NewUserService(repo)
	authSvc := service.NewAuthService(repo, jwtSecret)
	userHandler := handler.NewUserHandler(userSvc)
	authHandler := handler.NewAuthHandler(authSvc)

	r := chi.NewRouter()

	r.Use(middleware.Recover)
	r.Use(middleware.Logging)

	r.Post("/auth/signup", authHandler.Signup)
	r.Post("/auth/login", authHandler.Login)

	r.Route("/users", func(r chi.Router) {

		r.Use(middleware.Auth(jwtSecret))
		r.Get("/", userHandler.GetAllUsers)
		r.Get("/{id}", userHandler.GetUserByID)
		r.Put("/{id}", userHandler.UpdateUser)
		r.Delete("/{id}", userHandler.DeleteUser)
	})

	log.Println("Server starting on :8081")
	newServ := http.Server{Addr: ":8081", Handler: r}

	go newServ.ListenAndServe()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	<-quit

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := newServ.Shutdown(ctx); err != nil {
		log.Fatal("Server forced to shutdown", err)
	}

	log.Println("Server exited")

}
