package main

import (
	"net/http"
)

type RouterConfig struct {
	APIKey       string
	RealIPHeader string
}

func NewRouter(config RouterConfig) http.Handler {
	reqInfo := requestInfo{realIPHeader: config.RealIPHeader}
	resWriter := responseWriter{}

	authRequests := authMiddleware(config.APIKey, resWriter)
	logRequests := logMiddleware(reqInfo)

	mux := http.NewServeMux()
	mux.Handle("GET /health", healthHandler(resWriter))
	mux.Handle("GET /v1/check", authRequests(checkHandler(reqInfo, resWriter)))
	mux.Handle("/", notFoundHandler(resWriter))

	return logRequests(mux)
}
