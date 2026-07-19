// Package middleware provides HTTP middleware for cross-cutting concerns.
// rate_limit.go implements per-IP rate limiting using a token bucket algorithm.
// Each IP gets a bucket of tokens; each request consumes one. When empty,
// requests are rejected with 429 Too Many Requests until tokens refill.
package middleware

import (
	"net"
	"net/http"
	"sync"
	"time"

	"golang.org/x/time/rate"
)

// rateLimiter tracks per-IP rate limiters. Unexported because it's an
// internal implementation detail — only RateLimit() is public.
// RWMutex allows concurrent IP lookups (reads) while serializing new
// IP registrations (writes).
type rateLimiter struct {
	visitors map[string]*rate.Limiter
	mut      sync.RWMutex
}

func newRateLimiter() *rateLimiter {

	var newRateLimit rateLimiter
	newRateLimit.visitors = make(map[string]*rate.Limiter, 0)
	return &newRateLimit

}

// getLimiter returns the rate limiter for a given IP, creating one if it
// doesn't exist. Uses the double-check pattern: read-lock to check, then
// write-lock to create. The second check after acquiring the write lock
// prevents duplicate creation when two goroutines race for the same new IP.
// Manual unlock (not defer) is required because we switch between lock types.
func (r *rateLimiter) getLimiter(IP string) *rate.Limiter {

	r.mut.RLock()
	val, ok := r.visitors[IP]
	if ok {
		r.mut.RUnlock()
		return val
	}
	r.mut.RUnlock()
	r.mut.Lock()
	val, ok = r.visitors[IP]
	if ok {
		r.mut.Unlock()
		return val
	}
	newRate := rate.NewLimiter(rate.Every(time.Second), 1000)
	r.visitors[IP] = newRate
	r.mut.Unlock()

	return newRate
}

// RateLimit returns rate limiting middleware. Uses a closure to capture a
// single rateLimiter instance shared across all requests. net.SplitHostPort
// strips the port from RemoteAddr so all requests from the same IP share
// one bucket regardless of source port.
func RateLimit() func(http.Handler) http.Handler {

	rateLim := newRateLimiter()

	return func(next http.Handler) http.Handler {

		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

			ip, _, _ := net.SplitHostPort(r.RemoteAddr)

			newRate := rateLim.getLimiter(ip)
			if !newRate.Allow() {
				http.Error(w, "rate limit exceeded", http.StatusTooManyRequests)
				return
			}
			next.ServeHTTP(w, r)

		})
	}

}
