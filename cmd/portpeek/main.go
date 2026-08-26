package main

import (
	"fmt"
	"net/http"
	"os"
	"time"
)

func main() {
	logger := newLogger(os.Getenv("LOG_FORMAT"))

	cfg, err := loadConfig()
	if err != nil {
		logger.Error("failed to load config", "error", err)
		os.Exit(1)
	}

	router := NewRouter(RouterConfig{
		APIKey:       cfg.apiKey,
		RealIPHeader: cfg.realIPHeader,
		Logger:       logger,
	})

	server := &http.Server{
		Addr:              fmt.Sprintf(":%s", cfg.port),
		Handler:           router,
		ReadHeaderTimeout: 5 * time.Second,
		ReadTimeout:       10 * time.Second,
		WriteTimeout:      10 * time.Second,
		IdleTimeout:       60 * time.Second,
	}

	logger.Info("portpeek listening", "port", cfg.port)

	if err := server.ListenAndServe(); err != nil {
		logger.Error("server failed", "error", err)
		os.Exit(1)
	}
}
