package middleware

import "net/http"

func Authorize(requiredRole string) func(http.Handler) http.Handler {

	return func(next http.Handler) http.Handler {

		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

			role, ok := r.Context().Value(RoleKey).(string)
			if !ok {
				http.Error(w, "invalid role", http.StatusForbidden)
				return
			}

			if role != requiredRole {
				http.Error(w, "forbidden", http.StatusForbidden)
				return
			}

			next.ServeHTTP(w, r)

		})
	}
}
