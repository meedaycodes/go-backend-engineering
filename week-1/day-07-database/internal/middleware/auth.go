package middleware

import (
	"log"
	"net/http"
)

// Authorization is the expected bearer token for authenticating requests.
const Authorization = "Bearer secret-token"

// Auth checks the Authorization header against the expected token.
// Requests with a missing or invalid token receive a 401 Unauthorized response
// and are not passed to downstream handlers.
func Auth(next http.Handler) http.Handler {

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		authHeader := r.Header.Get("Authorization")

		if authHeader != Authorization {
			log.Println("Wrong Authorization credentials")
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}
		next.ServeHTTP(w, r)

	})
}
