package main

import (
	"log/slog"
	"net"
	"net/http"
	"strconv"
	"strings"
	"time"
)

func healthHandler(resWriter responseWriter) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		resWriter.text(w, http.StatusOK, "OK")
	})
}

func notFoundHandler(resWriter responseWriter) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		resWriter.text(w, http.StatusNotFound, "NOT_FOUND")
	})
}

func checkHandler(logger *slog.Logger, reqInfo requestInfo, resWriter responseWriter) http.Handler {
	dialer := net.Dialer{
		Timeout: 3 * time.Second,
	}

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		host := reqInfo.clientIP(r)
		port := strings.TrimSpace(r.URL.Query().Get("port"))

		if port == "" {
			resWriter.text(w, http.StatusBadRequest, "PORT_REQUIRED")
			return
		}

		portInt, err := strconv.Atoi(port)
		if err != nil || portInt < 1 || portInt > 65535 {
			resWriter.text(w, http.StatusBadRequest, "INVALID_PORT")
			return
		}

		address := net.JoinHostPort(host, port)

		conn, err := dialer.DialContext(r.Context(), "tcp", address)
		if err != nil {
			logger.Info("failed to dial", "address", address, "error", err)
			resWriter.text(w, http.StatusOK, "CLOSED")
			return
		}
		defer conn.Close()

		resWriter.text(w, http.StatusOK, "OPEN")
	})
}
