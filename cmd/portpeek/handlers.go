package main

import (
	"context"
	"net"
	"net/http"
	"strconv"
	"strings"
	"time"
)

func healthHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		writeText(w, http.StatusOK, "OK")
	}
}

func checkHandler(realIPHeader string, dialer *net.Dialer) http.HandlerFunc {
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

		ctx, cancel := context.WithTimeout(r.Context(), 3*time.Second)
		defer cancel()

		address := net.JoinHostPort(host, port)

		conn, err := dialer.DialContext(ctx, "tcp", address)
		if err != nil {
			writeText(w, http.StatusOK, "CLOSED")
			return
		}
		defer conn.Close()

		writeText(w, http.StatusOK, "OPEN")
	})
}
