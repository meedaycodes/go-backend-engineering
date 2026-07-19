// Package middleware provides HTTP middleware for cross-cutting concerns.
// auth.go validates JWT tokens on protected routes. It uses a closure pattern
// to inject the JWT secret into middleware with Chi's fixed signature.
package middleware

import (
	"context"
	"errors"
	"net/http"
	"strings"

	"github.com/golang-jwt/jwt/v5"
)

// contextKey is a custom type for context keys to prevent collisions.
// Using a bare string risks another package accidentally using the same key.
type contextKey string

// UserIDKey is the context key for storing the authenticated user's ID.
// Handlers read it with r.Context().Value(middleware.UserIDKey).
const UserIDKey contextKey = "userID"
const RoleKey contextKey = "role"

// Auth returns JWT validation middleware. Uses a closure (function returning
// a function) to capture the jwtSecret while conforming to Chi's middleware
// signature: func(http.Handler) http.Handler.
//
// Per request, it: extracts Bearer token from Authorization header, parses
// and verifies the JWT signature (HMAC-SHA256), checks expiry, extracts the
// user ID from the "sub" claim, and injects it into the request context.
// Returns 401 Unauthorized if any step fails.
func Auth(jwtSecret string) func(http.Handler) http.Handler {

	return func(next http.Handler) http.Handler {

		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

			authHeader := r.Header.Get("Authorization")

			if !strings.HasPrefix(authHeader, "Bearer ") {
				http.Error(w, "Invalid credentials", http.StatusUnauthorized)
				return
			}

			token := authHeader[7:]

			sToken, err := jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {

				_, ok := token.Method.(*jwt.SigningMethodHMAC)
				if !ok {
					return nil, errors.New("wrong signing Method")
				}
				return []byte(jwtSecret), nil
			})

			if err != nil {
				http.Error(w, "Invalid credentials", http.StatusUnauthorized)
				return
			}

			if !sToken.Valid {
				http.Error(w, "Invalid credentials", http.StatusUnauthorized)
				return
			}

			claims, ok := sToken.Claims.(jwt.MapClaims)
			if !ok {
				http.Error(w, "invalid token claims", http.StatusUnauthorized)
				return
			}

			userID, err := claims.GetSubject()
			if err != nil {
				http.Error(w, "Invalid credentials", http.StatusUnauthorized)
				return
			}
			role, ok := claims["role"].(string)
			if !ok {
				http.Error(w, "invalid token claims", http.StatusUnauthorized)
				return
			}
			ctx := context.WithValue(r.Context(), UserIDKey, userID)
			ctx = context.WithValue(ctx, RoleKey, role)

			next.ServeHTTP(w, r.WithContext(ctx))

		})

	}

}
