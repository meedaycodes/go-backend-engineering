// Package main is the entry point that wires all layers together and starts
// the HTTP server. Auth routes (/auth/signup, /auth/login) are public.
// User routes (/users) are protected behind JWT middleware.
// Dependencies flow: pool → repo → cache → worker → services → handlers → router → server.
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

	"github.com/meedaycodes/day14-integration-testing/internal/cache"
	"github.com/meedaycodes/day14-integration-testing/internal/config"
	"github.com/meedaycodes/day14-integration-testing/internal/handler"
	"github.com/meedaycodes/day14-integration-testing/internal/middleware"
	"github.com/meedaycodes/day14-integration-testing/internal/repository"
	"github.com/meedaycodes/day14-integration-testing/internal/service"
	"github.com/meedaycodes/day14-integration-testing/internal/worker"
	"github.com/prometheus/client_golang/prometheus/promhttp"
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

	// ctx is the application-lifetime context. Cancelling it (via defer cancel)
	// signals long-running goroutines — like the email worker — to shut down.
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Start the email worker as a background goroutine. It blocks on the jobs
	// channel and exits when ctx is cancelled at shutdown.
	emailWorker := worker.NewEmailWorker(100)
	go emailWorker.Start(ctx)

	readPool, err := pgxpool.New(context.Background(), cfg.DatabaseURL)
	if err != nil {
		slog.Error("Failed to connect to read replica", "error", err)
		os.Exit(1)
	}
	defer readPool.Close()

	repo := repository.NewReadWriteUserRepository(pool, readPool)
	userCache := cache.NewRedisCache(cfg.RedisAddr)
	userSvc := service.NewUserService(repo, userCache)
	authSvc := service.NewAuthService(repo, cfg.JWTSecret, emailWorker)
	userHandler := handler.NewUserHandler(userSvc)
	authHandler := handler.NewAuthHandler(authSvc)

	r := chi.NewRouter()

	r.Use(middleware.RateLimit())
	r.Use(middleware.Recover)
	r.Use(middleware.Logging)

	r.Get("/healthz", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	r.Get("/readyz", func(w http.ResponseWriter, r *http.Request) {
		if err := pool.Ping(r.Context()); err != nil {
			http.Error(w, "database unavailable", http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusOK)
	})

	r.Handle("/metrics", promhttp.Handler())
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

	// Run the server in a goroutine so main can block on the quit signal below.
	go func() {
		if err := newServ.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			slog.Error(err.Error())
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	<-quit

	// Give in-flight requests up to 10 seconds to complete before forcing exit.
	// shutdownctx is separate from the worker ctx — they serve different purposes.
	shutdownctx, shutdowncancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer shutdowncancel()

	if err := newServ.Shutdown(shutdownctx); err != nil {
		slog.Error("Server forced to shutdown", "error", err)
		os.Exit(1)
	}

	slog.Info("Server exited")

}
