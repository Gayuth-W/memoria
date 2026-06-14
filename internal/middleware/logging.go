package middleware

// This middleware provides request-level observability

import (
	"context"
	"log/slog"
	"net/http"
	"time"

	"github.com/google/uuid"
)

const RequestIDKey contextKey = "request_id"

// It wraps the ResponseWriter to capture HTTP status codes and response sizes
type statusRecorder struct {
	http.ResponseWriter
	status int
	bytes  int
}

func (s *statusRecorder) WriteHeader(code int) {
	s.status = code
	s.ResponseWriter.WriteHeader(code)
}

func (s *statusRecorder) Write(b []byte) (int, error) {
	n, err := s.ResponseWriter.Write(b)
	s.bytes += n
	return n, err
}

func RequestLogger(logger *slog.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// measures request latency
			start := time.Now()

			// generates a unique request ID for traceability
			reqID := uuid.New().String()

			// Make the request ID available to downstream handlers.
			ctx := context.WithValue(r.Context(), RequestIDKey, reqID)

			// Wrap the ResponseWriter so we can capture the final status code and response size for logging.
			rec := &statusRecorder{ResponseWriter: w, status: 200}

			// Execute the next handler with the request ID attached to the request context.
			next.ServeHTTP(rec, r.WithContext(ctx))

			// logs structured request metadata using slog after the request completes
			logger.Info("request",
				slog.String("request_id", reqID),
				slog.String("method", r.Method),
				slog.String("path", r.URL.Path),
				slog.Int("status", rec.status),
				slog.Int("bytes", rec.bytes),
				slog.Duration("latency", time.Since(start)),
			)
		})
	}
}
