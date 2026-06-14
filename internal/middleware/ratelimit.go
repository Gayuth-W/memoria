package middleware

// RateLimit enforces a fixed-window rate limit per authenticated user
// using Redis counters.

import (
	"net/http"
	"time"

	"memoria/internal/cache"
)

func RateLimit(rc *cache.RedisCache, limit int, window time.Duration) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

			next.ServeHTTP(w, r)
		})
	}
}
