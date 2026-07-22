package main

import (
	"errors"
	"log"
	"net/http"
)

func main() {
	config, err := loadConfig()
	if err != nil {
		log.Fatalf("failed to load config: %v", err)
	}

	httpServer := newServer(config)

	log.Printf("PortPeek listening on %s", config.port)
	if err := httpServer.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
		log.Fatalf("server failed: %v", err)
	}
}
