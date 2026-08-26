package main

import (
	"log/slog"
	"net"
	"net/http"
	"strconv"
	"strings"
	"time"
)

func healthHandler() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		writeText(w, http.StatusOK, "OK")
	})
}

func notFoundHandler() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		writeText(w, http.StatusNotFound, "NOT_FOUND")
	})
}

func checkHandler(realIPHeader string) http.Handler {
	dialer := net.Dialer{
		Timeout: 3 * time.Second,
	}

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		host := getClientIP(r, realIPHeader)
		port := strings.TrimSpace(r.URL.Query().Get("port"))

		if port == "" {
			writeText(w, http.StatusBadRequest, "PORT_REQUIRED")
			return
		}

		portInt, err := strconv.Atoi(port)
		if err != nil || portInt < 1 || portInt > 65535 {
			writeText(w, http.StatusBadRequest, "INVALID_PORT")
			return
		}

		address := net.JoinHostPort(host, port)

		conn, err := dialer.DialContext(r.Context(), "tcp", address)
		if err != nil {
			slog.Info("failed to dial", "address", address, "error", err)
			writeText(w, http.StatusOK, "CLOSED")
			return
		}
		defer conn.Close()

		writeText(w, http.StatusOK, "OPEN")
	})
}
