package main

import (
	"net"
	"net/http"
)

type RouterConfig struct {
	APIKey       string
	RealIPHeader string
}

func NewRouter(config RouterConfig) http.Handler {
	authRequests := authMiddleware(config.APIKey)
	logRequests := logMiddleware(config.RealIPHeader)

	mux := http.NewServeMux()
	mux.Handle("GET /health", healthHandler())
	mux.Handle("GET /v1/check", authRequests(checkHandler(config.RealIPHeader, &net.Dialer{})))
	mux.Handle("/", notFoundHandler())

	return logRequests(mux)
}
