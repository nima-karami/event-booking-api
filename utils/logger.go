package utils

import (
	"log/slog"
	"os"
)

var Logger *slog.Logger

func InitLogger() {
	env := GetEnvString("ENV", "development")

	var handler slog.Handler

	if env == "production" {
		// JSON format for production (machine-readable, for log aggregation)
		handler = slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
			Level: slog.LevelInfo,
		})
	} else {
		// Human-readable format for development
		handler = slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
			Level: slog.LevelDebug,
		})
	}

	Logger = slog.New(handler)
	slog.SetDefault(Logger) // Set as default logger for entire application
}
