package main

import (
	"log/slog"
	"net/http"
)

type RouterConfig struct {
	APIKey       string
	RealIPHeader string
	Logger       *slog.Logger
}

func NewRouter(config RouterConfig) http.Handler {
	reqInfo := requestInfo{realIPHeader: config.RealIPHeader}
	resWriter := responseWriter{logger: config.Logger}

	authRequests := authMiddleware(config.APIKey, resWriter)
	logRequests := logMiddleware(config.Logger, reqInfo)

	mux := http.NewServeMux()
	mux.Handle("GET /health", healthHandler(resWriter))
	mux.Handle("GET /v1/check", authRequests(checkHandler(config.Logger, reqInfo, resWriter)))
	mux.Handle("/", notFoundHandler(resWriter))

	return logRequests(mux)
}
