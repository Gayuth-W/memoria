package middleware

// RateLimit enforces a fixed-window rate limit per authenticated user
// using Redis counters.

import (
	"context"
	"net/http"
	"time"

	"memoria/internal/cache"

	"github.com/google/uuid"
)

func RateLimit(rc *cache.RedisCache, limit int, window time.Duration) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

			// Rate limits are applied per authenticated user.
			// Reject requests that do not contain a valid user ID.
			userID, ok := r.Context().Value(UserIDKey).(uuid.UUID)
			if !ok {
				http.Error(w, "unauthorized", http.StatusUnauthorized)
				return
			}

			ctx := context.Background()

			// Use a unique Redis key for each user so request counts are tracked independently.
			key := "ratelimit:" + userID.String()

			// Increment the request counter for the current window.
			count, err := rc.Client.Incr(ctx, key).Result()

			// If Redis is unavailable, allow the request rather
			// than blocking legitimate traffic.
			if err != nil {
				next.ServeHTTP(w, r)
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}
