package middleware

import (
	"context"
	"net/http"

	"memoria/internal/repository"
)

type contextKey string

const UserIDKey contextKey = "user_id"

func APIKeyAuth(repo *repository.UserRepo) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

			key := r.Header.Get("X-API-Key")
			if key == "" {
				http.Error(w, "missing api key", 401)
				return
			}

			user, err := repo.GetByAPIKey(key)
			if err != nil {
				http.Error(w, "invalid api key", 401)
				return
			}

			ctx := context.WithValue(r.Context(), UserIDKey, user.ID)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
