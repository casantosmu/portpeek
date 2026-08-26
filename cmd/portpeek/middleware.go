package main

import (
	"crypto/subtle"
	"log/slog"
	"net/http"
	"time"
)

type responseWriter struct {
	http.ResponseWriter
	statusCode int
}

func (w *responseWriter) WriteHeader(statusCode int) {
	w.statusCode = statusCode
	w.ResponseWriter.WriteHeader(statusCode)
}

func (w *responseWriter) Write(data []byte) (int, error) {
	return w.ResponseWriter.Write(data)
}

type middleware func(http.Handler) http.Handler

func authMiddleware(apiKey string) middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			providedKey := r.Header.Get("X-API-Key")
			if subtle.ConstantTimeCompare([]byte(providedKey), []byte(apiKey)) != 1 {
				writeText(w, http.StatusUnauthorized, "UNAUTHORIZED")
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}

func logMiddleware(realIPHeader string) middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.URL.Path == "/health" {
				next.ServeHTTP(w, r)
				return
			}

			start := time.Now()

			rw := &responseWriter{
				ResponseWriter: w,
				statusCode:     http.StatusOK,
			}
			next.ServeHTTP(rw, r)

			slog.Info(
				"http request",
				"method", r.Method,
				"route", r.Pattern,
				"status_code", rw.statusCode,
				"duration_ms", time.Since(start).Milliseconds(),
				"client_ip", getClientIP(r, realIPHeader),
			)
		})
	}
}
