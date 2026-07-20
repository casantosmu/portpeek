package main

import (
	"context"
	"crypto/subtle"
	"errors"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"
)

func main() {
	apiKey := strings.TrimSpace(os.Getenv("API_KEY"))
	if apiKey == "" {
		log.Fatal("API_KEY environment variable is required")
	}

	port := strings.TrimSpace(os.Getenv("PORT"))
	if port == "" {
		port = "8080"
	}

	clientIPHeader := strings.TrimSpace(os.Getenv("CLIENT_IP_HEADER"))
	dialer := &net.Dialer{}

	mux := http.NewServeMux()
	mux.HandleFunc("GET /health", healthHandler())
	mux.HandleFunc("GET /v1/check", checkHandler(apiKey, clientIPHeader, dialer))

	httpServer := &http.Server{
		Addr:              fmt.Sprintf(":%s", port),
		Handler:           requestLogger(mux, clientIPHeader),
		ReadHeaderTimeout: 5 * time.Second,
		ReadTimeout:       10 * time.Second,
		WriteTimeout:      10 * time.Second,
		IdleTimeout:       60 * time.Second,
	}

	log.Printf("PortPeek listening on %s", port)
	if err := httpServer.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
		log.Fatalf("server failed: %v", err)
	}
}

func healthHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		writeText(w, http.StatusOK, "OK")
	}
}

func checkHandler(apiKey string, clientIPHeader string, dialer *net.Dialer) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		providedKey := r.Header.Get("X-API-Key")
		if subtle.ConstantTimeCompare([]byte(providedKey), []byte(apiKey)) != 1 {
			writeText(w, http.StatusUnauthorized, "UNAUTHORIZED")
			return
		}

		host := strings.TrimSpace(r.URL.Query().Get("host"))
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

		if host == "" {
			host = getClientIP(r, clientIPHeader)
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
	}
}

func writeText(w http.ResponseWriter, statusCode int, value string) {
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(statusCode)

	if _, err := fmt.Fprintln(w, value); err != nil {
		log.Printf("failed to write response: %v", err)
	}
}

func getClientIP(r *http.Request, header string) string {
	if header != "" {
		value := strings.TrimSpace(r.Header.Get(header))
		ip := net.ParseIP(value)

		if ip == nil {
			return ""
		}

		return ip.String()
	}

	remoteAddr := strings.TrimSpace(r.RemoteAddr)

	host, _, err := net.SplitHostPort(remoteAddr)
	if err != nil {
		host = remoteAddr
	}

	ip := net.ParseIP(host)
	if ip == nil {
		return ""
	}

	return ip.String()
}

func requestLogger(next http.Handler, clientIPHeader string) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		startedAt := time.Now()
		next.ServeHTTP(w, r)
		log.Printf(
			"method=%s path=%q duration_ms=%d ip=%q",
			r.Method, r.URL.Path, time.Since(startedAt).Milliseconds(), getClientIP(r, clientIPHeader),
		)
	})
}
