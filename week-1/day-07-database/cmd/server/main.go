// Package main is the entry point that wires all layers together:
// database pool → repository → service → handler → router → server.
// Each layer receives its dependency through constructor injection.
// The server runs with graceful shutdown — on SIGINT/SIGTERM, it stops
// accepting new connections and waits up to 10 seconds for active requests
// to complete before exiting.
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

	"github.com/meedaycodes/day-07-database/internal/handler"
	"github.com/meedaycodes/day-07-database/internal/middleware"
	"github.com/meedaycodes/day-07-database/internal/repository"
	"github.com/meedaycodes/day-07-database/internal/service"
)

func main() {

	dbURL := "postgres://habeebaramideshomuyiwa@localhost:5432/day07_users?sslmode=disable"

	pool, err := pgxpool.New(context.Background(), dbURL)

	if err != nil {
		log.Fatal(err)
	}
	log.Println("Connected to database")

	defer pool.Close()

	repo := repository.NewPostgresUserRepository(pool)
	svc := service.NewUserService(repo)
	servHandler := handler.NewUserHandler(svc)

	r := chi.NewRouter()

	r.Use(middleware.Recover)
	r.Use(middleware.Logging)

	r.Route("/users", func(r chi.Router) {

		r.Use(middleware.Auth)
		r.Post("/", servHandler.CreateUser)
		r.Get("/", servHandler.GetAllUsers)
		r.Get("/{id}", servHandler.GetUserByID)
		r.Put("/{id}", servHandler.UpdateUser)
		r.Delete("/{id}", servHandler.DeleteUser)
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
