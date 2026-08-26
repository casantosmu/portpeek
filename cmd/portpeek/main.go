package main

import (
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"time"
)

func main() {
	cfg, err := setup()
	if err != nil {
		slog.Error("failed to load config", "error", err)
		os.Exit(1)
	}

	router := NewRouter(RouterConfig{
		APIKey:       cfg.apiKey,
		RealIPHeader: cfg.realIPHeader,
	})

	server := &http.Server{
		Addr:              fmt.Sprintf(":%s", cfg.port),
		Handler:           router,
		ReadHeaderTimeout: 5 * time.Second,
		ReadTimeout:       10 * time.Second,
		WriteTimeout:      10 * time.Second,
		IdleTimeout:       60 * time.Second,
	}

	slog.Info("portpeek listening", "port", cfg.port)

	if err := server.ListenAndServe(); err != nil {
		slog.Error("server failed", "error", err)
		os.Exit(1)
	}
}
