package observability

import (
	"log/slog"
	"os"
)

// NewLogger creates a structured logger using Go's slog package.
func NewLogger() *slog.Logger {
	// It defaults to INFO level logging, but switches to DEBUG if the LOG_LEVEL environment variable is set to debug.
	level := slog.LevelInfo
	if os.Getenv("LOG_LEVEL") == "debug" {
		level = slog.LevelDebug
	}

	// The logger outputs logs in JSON format to stdout, making them suitable for log aggregation systems
	logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: level}))

	// registers the logger as the application's default logger using slog
	slog.SetDefault(logger)
	return logger
}
