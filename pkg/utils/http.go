package utils

import (
	"log"
	"net/http"
)

// GetHeader captures the HTTP heads
func GetHeader(key string, r *http.Request) string {
	v := r.Header.Get(key)
	log.Printf("Header[%s] %s", key, v)
	return v
}
