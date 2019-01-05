package main

import (
	"log"
	"net/http"
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

func getHeader(key string, r *http.Request) string {
	v := r.Header.Get(key)
	log.Printf("Header[%s] %s", key, v)
	return v
}
