package main

import (
	"fmt"
	"log/slog"
	"net/http"
)

type responseWriter struct {
	logger *slog.Logger
}

func (rw responseWriter) text(w http.ResponseWriter, statusCode int, value string) {
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(statusCode)

	if _, err := fmt.Fprintln(w, value); err != nil {
		rw.logger.Error("failed to write response", "error", err)
	}
}
