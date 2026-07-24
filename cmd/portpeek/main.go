package main

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"time"
)

func main() {
	cfg, err := loadConfig()
	if err != nil {
		log.Fatalf("failed to load config: %v", err)
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

	log.Printf("PortPeek listening on %s", cfg.port)

	if err := server.ListenAndServe(); err != nil &&
		!errors.Is(err, http.ErrServerClosed) {
		log.Fatalf("server failed: %v", err)
	}
}
