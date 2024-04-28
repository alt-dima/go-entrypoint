package utils

import (
	"log/slog"
	"os"
)

var Logger = slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
	Level: slog.LevelDebug,
	ReplaceAttr: func(groups []string, a slog.Attr) slog.Attr {
		if a.Key == slog.TimeKey {
			return slog.Attr{Key: "timestamp", Value: a.Value}
		}
		if a.Key == slog.MessageKey {
			return slog.Attr{Key: "message", Value: a.Value}
		}
		return a
	},
}))
