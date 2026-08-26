package main

import (
	"fmt"
	"log"
	"net/http"
)

func writeText(w http.ResponseWriter, statusCode int, value string) {
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(statusCode)

	if _, err := fmt.Fprintln(w, value); err != nil {
		log.Printf("failed to write response: %v", err)
	}
}
