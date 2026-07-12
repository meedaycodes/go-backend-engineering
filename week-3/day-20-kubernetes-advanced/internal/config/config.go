// Package config loads application configuration from environment variables.
// In development, a .env file is read by godotenv. In production, env vars
// are set directly by Docker/Kubernetes. Required values (DatabaseURL, JWTSecret)
// cause a startup failure if missing — fail fast rather than crash later.
package config

import (
	"errors"
	"os"

	"log/slog"

	"github.com/joho/godotenv"
)

// Config holds all application configuration. Struct fields map to environment
// variables: DATABASE_URL, JWT_SECRET, SERVER_PORT, LOG_LEVEL.
// Using a struct (not individual globals) keeps config grouped, testable,
// and injectable as a single dependency.
type Config struct {
	DatabaseURL,
	JWTSecret,
	ServerPort,
	LogLevel string
	RedisAddr string
}

// Load reads configuration from environment variables. godotenv.Load() reads
// .env for local development — its failure is ignored because production
// environments set env vars directly. Required fields (DatabaseURL, JWTSecret)
// return an error if empty. Optional fields (ServerPort, LogLevel) fall back
// to sensible defaults.
func Load() (conf *Config, err error) {

	err = godotenv.Load()
	if err != nil {
		slog.Error(err.Error())
	}

	databaseURL := os.Getenv("DATABASE_URL")
	if databaseURL == "" {
		return conf, errors.New("invalid database credentials")
	}

	jwtSecret := os.Getenv("JWT_SECRET")
	if jwtSecret == "" {
		return conf, errors.New("invalid secret cannot be empty")
	}

	serverPort := os.Getenv("SERVER_PORT")
	if serverPort == "" {
		serverPort = "8081"
	}

	logLevel := os.Getenv("LOG_LEVEL")
	if logLevel == "" {
		logLevel = "info"
	}

	redisAddr := os.Getenv("REDIS_ADDR")
	if redisAddr == "" {
		redisAddr = "localhost:6379"
	}

	conf = &Config{DatabaseURL: databaseURL, JWTSecret: jwtSecret, ServerPort: serverPort, LogLevel: logLevel, RedisAddr: redisAddr}

	return conf, nil

}

// SetupLogger creates a structured JSON logger with the given minimum log level.
// Uses Go's built-in slog package with a JSONHandler writing to stdout.
// The logger is set as the global default via slog.SetDefault, so any call to
// slog.Info/Error/etc anywhere in the app uses this configuration.
// Returns the logger for explicit injection where needed.
func SetupLogger(logLevel string) *slog.Logger {

	var slogLevel slog.Level

	switch logLevel {
	case "debug":
		slogLevel = slog.LevelDebug
	case "info":
		slogLevel = slog.LevelInfo
	case "warn":
		slogLevel = slog.LevelWarn
	case "error":
		slogLevel = slog.LevelError
	default:
		slogLevel = slog.LevelInfo
	}

	options := slog.HandlerOptions{Level: slogLevel}
	handlerLogger := slog.NewJSONHandler(os.Stdout, &options)

	logger := slog.New(handlerLogger)
	slog.SetDefault(logger)

	return logger
}
