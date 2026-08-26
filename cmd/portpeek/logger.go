package main

import (
	"log/slog"
	"os"
	"strings"
)

func newLogger(format string) *slog.Logger {
	if strings.EqualFold(strings.TrimSpace(format), "json") {
		return slog.New(slog.NewJSONHandler(os.Stderr, nil))
	}

	return slog.New(slog.NewTextHandler(os.Stderr, nil))
}
