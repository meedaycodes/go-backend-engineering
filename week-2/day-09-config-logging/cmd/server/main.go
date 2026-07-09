// Package main is the entry point that wires all layers together and starts
// the HTTP server. Auth routes (/auth/signup, /auth/login) are public.
// User routes (/users) are protected behind JWT middleware.
// Dependencies flow: pool → repo → services → handlers → router → server.
package main

import (
	"context"
	"log"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/meedaycodes/day09-config-logging/internal/config"
	"github.com/meedaycodes/day09-config-logging/internal/handler"
	"github.com/meedaycodes/day09-config-logging/internal/middleware"
	"github.com/meedaycodes/day09-config-logging/internal/repository"
	"github.com/meedaycodes/day09-config-logging/internal/service"
)

func main() {

	cfg, err := config.Load()

	if err != nil {
		log.Fatal(err)
	}

	_ = config.SetupLogger(cfg.LogLevel)

	pool, err := pgxpool.New(context.Background(), cfg.DatabaseURL)

	if err != nil {
		slog.Error("Failed to connect to database", "error", err)
		os.Exit(1)
	}
	slog.Info("Connected to database")

	defer pool.Close()

	repo := repository.NewPostgresUserRepository(pool)
	userSvc := service.NewUserService(repo)
	authSvc := service.NewAuthService(repo, cfg.JWTSecret)
	userHandler := handler.NewUserHandler(userSvc)
	authHandler := handler.NewAuthHandler(authSvc)

	r := chi.NewRouter()

	r.Use(middleware.Recover)
	r.Use(middleware.Logging)

	r.Post("/auth/signup", authHandler.Signup)
	r.Post("/auth/login", authHandler.Login)

	r.Route("/users", func(r chi.Router) {

		r.Use(middleware.Auth(cfg.JWTSecret))
		r.Get("/", userHandler.GetAllUsers)
		r.Get("/{id}", userHandler.GetUserByID)
		r.Put("/{id}", userHandler.UpdateUser)
		r.Delete("/{id}", userHandler.DeleteUser)
	})

	slog.Info("Server starting", "port", cfg.ServerPort)
	newServ := http.Server{Addr: ":" + cfg.ServerPort, Handler: r}

	go newServ.ListenAndServe()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	<-quit

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := newServ.Shutdown(ctx); err != nil {
		slog.Error("Server forced to shutdown", "error", err)
		os.Exit(1)
	}

	slog.Info("Server exited")

}
