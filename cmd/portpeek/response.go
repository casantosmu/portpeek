package main

import (
	"fmt"
	"log"
	"net/http"
)

func writeText(w http.ResponseWriter, r *http.Request, statusCode int, value string) {
	if r.URL.Path != "/health" {
		log.Printf("status=%d", statusCode)
	}

	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(statusCode)

	if _, err := fmt.Fprintln(w, value); err != nil {
		log.Printf("failed to write response: %v", err)
	}
}
