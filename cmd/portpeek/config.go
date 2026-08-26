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
	var missing []string

	apiKey := strings.TrimSpace(os.Getenv("API_KEY"))
	if apiKey == "" {
		missing = append(missing, "API_KEY")
	}

	realIPHeader := strings.TrimSpace(os.Getenv("REAL_IP_HEADER"))
	if realIPHeader == "" {
		missing = append(missing, "REAL_IP_HEADER")
	}

	if len(missing) > 0 {
		return config{}, errors.New(
			"missing required environment variable(s): " + strings.Join(missing, ", "),
		)
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
