package logger

import (
	"log/slog"
	"os"
)

type Logger struct {
	*slog.Logger
}

func New(env string) *Logger {
	level := slog.LevelInfo
	if env == "local" {
		level = slog.LevelDebug
	}

	handler := slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level: level,
	})

	return &Logger{
		Logger: slog.New(handler),
	}
}
