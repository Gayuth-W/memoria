package middleware

// This middleware provides request-level observability

import (
	"net/http"
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
