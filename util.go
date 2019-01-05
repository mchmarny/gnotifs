package main

import (
	"log"
	"os"
)

func mustGetEnv(key, fallbackValue string) string {
	if val, ok := os.LookupEnv(key); ok {
		return val
	}

	if fallbackValue == "" {
		log.Fatalf("Required env var (%s) not set", key)
	}

	return fallbackValue
}
