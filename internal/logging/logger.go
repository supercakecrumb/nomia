package logging

import (
	"log/slog"
	"os"
)

// InitLogger initializes the global logger with the specified level and format
func InitLogger(level string, format string) *slog.Logger {
	// Parse log level
	var logLevel slog.Level
	switch level {
	case "debug":
		logLevel = slog.LevelDebug
	case "info":
		logLevel = slog.LevelInfo
	case "warn":
		logLevel = slog.LevelWarn
	case "error":
		logLevel = slog.LevelError
	default:
		logLevel = slog.LevelInfo
	}

	// Create handler options
	opts := &slog.HandlerOptions{
		Level:     logLevel,
		AddSource: true, // Add source file and line number
	}

	// Create handler based on format
	var handler slog.Handler
	if format == "text" {
		handler = slog.NewTextHandler(os.Stdout, opts)
	} else {
		// JSON format is the default
		handler = slog.NewJSONHandler(os.Stdout, opts)
	}

	// Create and set default logger
	logger := slog.New(handler)
	slog.SetDefault(logger)

	return logger
}
