package middleware

import (
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/meedaycodes/day25-performance/internal/metrics"
)

// wrappedResponseWriter wraps http.ResponseWriter to capture the status code
// written by downstream handlers, since the standard ResponseWriter does not expose it.
type wrappedResponseWriter struct {
	http.ResponseWriter
	statusCode int
}

// WriteHeader captures the status code and delegates to the underlying ResponseWriter.
func (w *wrappedResponseWriter) WriteHeader(code int) {
	w.statusCode = code
	w.ResponseWriter.WriteHeader(code)
}

// Logging logs the HTTP method, path, status code, and duration of each request.
func Logging(next http.Handler) http.Handler {

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		start := time.Now()
		wrapped := &wrappedResponseWriter{ResponseWriter: w, statusCode: http.StatusOK}

		next.ServeHTTP(wrapped, r)

		log.Printf("%s %s %d %s", r.Method, r.URL.Path, wrapped.statusCode, time.Since(start))

		strStatus := strconv.Itoa(wrapped.statusCode)
		metrics.RequestsTotal.WithLabelValues(r.Method, r.URL.Path, strStatus).Inc()
		metrics.RequestDuration.WithLabelValues(r.Method, r.URL.Path).Observe(time.Since(start).Seconds())
	})
}
