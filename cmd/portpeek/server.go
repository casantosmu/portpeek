package main

import (
	"fmt"
	"net"
	"net/http"
	"time"
)

func newServer(config config) *http.Server {
	authRequests := authMiddleware(config.apiKey)
	logRequests := logMiddleware(config.realIPHeader)

	mux := http.NewServeMux()
	mux.Handle("GET /health", healthHandler())
	mux.Handle("GET /v1/check", authRequests(checkHandler(config.realIPHeader, &net.Dialer{})))

	return &http.Server{
		Addr:              fmt.Sprintf(":%s", config.port),
		Handler:           logRequests(mux),
		ReadHeaderTimeout: 5 * time.Second,
		ReadTimeout:       10 * time.Second,
		WriteTimeout:      10 * time.Second,
		IdleTimeout:       60 * time.Second,
	}
}
