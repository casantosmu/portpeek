package main

import (
	"net/http"
	"strings"
)

type requestInfo struct {
	realIPHeader string
}

func (ri requestInfo) clientIP(r *http.Request) string {
	return strings.TrimSpace(r.Header.Get(ri.realIPHeader))
}
