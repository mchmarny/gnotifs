package utils

import (
	"log"
	"os"
)

// MustGetEnv looks up env vars and errs if now fallback value
func MustGetEnv(key, fallbackValue string) string {
	if val, ok := os.LookupEnv(key); ok {
		return val
	}

	if fallbackValue == "" {
		log.Fatalf("Required env var (%s) not set", key)
	}

	return fallbackValue
}
