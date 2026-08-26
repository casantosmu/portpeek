package main

import (
	"errors"
	"log/slog"
	"os"
	"strings"
)

type config struct {
	apiKey       string
	port         string
	realIPHeader string
}

func setup() (config, error) {
	if strings.EqualFold(strings.TrimSpace(os.Getenv("LOG_FORMAT")), "json") {
		slog.SetDefault(slog.New(slog.NewJSONHandler(os.Stderr, nil)))
	} else {
		slog.SetDefault(slog.New(slog.NewTextHandler(os.Stderr, nil)))
	}

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
