package main

import (
	"crypto/subtle"
	"log"
	"net/http"
	"time"
)

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
			startedAt := time.Now()
			next.ServeHTTP(w, r)
			log.Printf(
				"method=%s path=%q duration_ms=%d ip=%q",
				r.Method, r.URL.Path, time.Since(startedAt).Milliseconds(), getClientIP(r, realIPHeader),
			)
		})
	}
}
