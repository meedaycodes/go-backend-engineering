// Package main is the entry point for the user API server.
// It wires together all layers (repository, service, handler), applies middleware,
// registers routes with the chi router, and starts the server with graceful shutdown.
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

	"github.com/meedaycodes/day-06-http-middleware/internal/handler"
	"github.com/meedaycodes/day-06-http-middleware/internal/middleware"
	"github.com/meedaycodes/day-06-http-middleware/internal/repository"
	"github.com/meedaycodes/day-06-http-middleware/internal/service"
)

func main() {

	repo := repository.NewInMemoryUserRepository()
	serv := service.NewUserService(repo)
	servHandler := handler.NewUserHandler(serv)

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

	go func() {
		if err := newServ.Shutdown(ctx); err != nil {
			log.Fatal("Server foced to shutdown", err)
		}

	}()

	log.Println("Server exited")

}
