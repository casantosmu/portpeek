package main

import (
	"errors"
	"os"
	"strings"
)

type config struct {
	apiKey       string
	port         string
	realIPHeader string
}

func loadConfig() (config, error) {
	apiKey := strings.TrimSpace(os.Getenv("API_KEY"))
	if apiKey == "" {
		return config{}, errors.New("API_KEY environment variable is required")
	}

	realIPHeader := strings.TrimSpace(os.Getenv("REAL_IP_HEADER"))
	if realIPHeader == "" {
		return config{}, errors.New("REAL_IP_HEADER environment variable is required")
	}

	port := strings.TrimSpace(os.Getenv("PORT"))
	if port == "" {
		port = "8080"
	}

	return config{
		apiKey:       apiKey,
		port:         port,
		realIPHeader: realIPHeader,
	}, nil
}
