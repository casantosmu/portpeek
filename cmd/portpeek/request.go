package main

import (
	"net/http"
	"strings"
)

func getClientIP(r *http.Request, header string) string {
	return strings.TrimSpace(r.Header.Get(header))
}
