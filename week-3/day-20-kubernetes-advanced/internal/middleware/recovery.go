// Package middleware provides HTTP middleware functions for cross-cutting concerns
// such as logging, authentication, and panic recovery.
package middleware

import (
	"log"
	"net/http"
)

// Recover catches panics from downstream handlers, logs the error,
// and returns a 500 Internal Server Error instead of crashing the server.
func Recover(next http.Handler) http.Handler {

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		defer func() {
			if err := recover(); err != nil {
				log.Println("recovered from panic", err)
				w.WriteHeader(http.StatusInternalServerError)
			}
		}()

		next.ServeHTTP(w, r)

	})
}
