// Package service contains business logic for the application.
// auth_service.go handles authentication operations: signup, login, and JWT
// token generation. Passwords are hashed with bcrypt (never stored in plain text).
// JWTs carry the user ID and are signed with HMAC-SHA256.
package service

import (
	"context"
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"

	"github.com/meedaycodes/day14-integration-testing/internal/model"
	"github.com/meedaycodes/day14-integration-testing/internal/repository"
	"github.com/meedaycodes/day14-integration-testing/internal/worker"
)

// Sentinel errors for authentication failures. Exported so handlers can
// match against them with errors.Is() and return appropriate HTTP status codes.
// ErrInvalidCredentials is deliberately vague — never reveal whether the email
// or password was wrong, to prevent user enumeration attacks.
var (
	ErrInvalidCredentials = errors.New("wrong email or password provided")
	ErrEmailAlreadyExists = errors.New("there is existing user with this email")
)

// AuthService handles signup, login, and token generation.
// It depends on UserRepository (interface) for data access, holds the JWT
// secret for signing tokens, and holds an EmailWorker for dispatching
// post-signup notifications. All three are injected via the constructor.
type AuthService struct {
	repo      repository.UserRepository
	jwtSecret string
	worker    *worker.EmailWorker
}

// NewAuthService creates a new AuthService with the given repository, JWT
// secret, and email worker. The worker is used to dispatch a welcome email
// after a successful signup — fire and forget, signup does not fail if the
// worker is slow or the channel is temporarily full.
func NewAuthService(repo repository.UserRepository, jwtSecret string, worker *worker.EmailWorker) *AuthService {

	newAuthService := &AuthService{
		repo:      repo,
		jwtSecret: jwtSecret,
		worker:    worker,
	}

	return newAuthService
}

// Signup registers a new user. It checks that the email is not already taken,
// hashes the password with bcrypt, saves the user, dispatches a welcome email
// job to the worker (non-blocking), and returns access + refresh tokens.
// Bcrypt uses a random salt so the same password produces different hashes
// each time — verified with CompareHashAndPassword, not string comparison.
func (a *AuthService) Signup(ctx context.Context, cu model.CreateUserRequest) (res model.AuthResponse, err error) {

	_, err = a.repo.FindByEmail(ctx, cu.Email)

	if err == nil {
		return res, ErrEmailAlreadyExists
	}

	if err != repository.ErrUserNotFound {
		return res, err
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(cu.Password), bcrypt.DefaultCost)
	if err != nil {
		return res, err
	}

	id := uuid.New().String()
	user := model.User{ID: id, Name: cu.Name, Email: cu.Email, PasswordHash: string(hashedPassword)}

	saveErr := a.repo.Save(ctx, user)
	if saveErr != nil {
		return res, saveErr
	}

	a.worker.Send(worker.EmailJob{To: user.Email, Subject: "Welcome", Body: "Thanks for signing up"})

	accessToken, err := a.generateToken(user.ID, time.Minute*15)
	if err != nil {
		return res, err
	}
	refreshTokenString, err := a.generateToken(user.ID, time.Hour*7*24)
	if err != nil {
		return res, err
	}

	res = model.AuthResponse{AccessToken: accessToken, RefreshToken: refreshTokenString}

	return res, nil

}

// generateToken creates a signed JWT with the user ID as the subject claim.
// Unexported because it's an internal helper — called by Signup and Login.
// The duration parameter controls token lifetime: 15 minutes for access tokens,
// 7 days for refresh tokens. Signed with HMAC-SHA256 using the service's secret.
func (a *AuthService) generateToken(userID string, duration time.Duration) (string, error) {

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub": userID,
		"exp": time.Now().Add(duration).Unix(),
		"iat": time.Now().Unix(),
	})

	tokenString, err := token.SignedString([]byte(a.jwtSecret))
	if err != nil {
		return tokenString, err
	}

	return tokenString, nil

}

// Login authenticates an existing user. Finds the user by email, verifies the
// password against the stored bcrypt hash using CompareHashAndPassword, and
// returns access + refresh tokens on success. Returns ErrInvalidCredentials
// for both wrong email and wrong password — never reveals which part failed.
func (a *AuthService) Login(ctx context.Context, lg model.LoginRequest) (auth model.AuthResponse, err error) {

	user, err := a.repo.FindByEmail(ctx, lg.Email)

	if err == repository.ErrUserNotFound {
		return auth, ErrInvalidCredentials
	}

	if err != nil {
		return auth, err
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(lg.Password))
	if err != nil {
		return auth, ErrInvalidCredentials
	}

	accessToken, err := a.generateToken(user.ID, time.Minute*15)
	if err != nil {
		return auth, err
	}
	refreshTokenString, err := a.generateToken(user.ID, time.Hour*7*24)
	if err != nil {
		return auth, err
	}

	auth = model.AuthResponse{AccessToken: accessToken, RefreshToken: refreshTokenString}

	return auth, nil

}
