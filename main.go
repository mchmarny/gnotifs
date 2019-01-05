package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

var (
	knownPublisherToken = ""
)

func defaultHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	msg := struct {
		Handlers []string `json:"handlers"`
	}{
		[]string{"POST: /gcs"},
	}
	json.NewEncoder(w).Encode(msg)
}

func main() {

	log.Print("Starting server...")

	// handlers
	http.HandleFunc("/", defaultHandler)
	http.HandleFunc("/gcs", gcsObjectHandler)
	http.HandleFunc("/_healthz", func(w http.ResponseWriter, _ *http.Request) {
		fmt.Fprint(w, "ok")
	})

	// variables
	knownPublisherToken = mustGetEnv("KNOWN_PUBLISHER_TOKEN", "")
	port := mustGetEnv("PORT", "8080")

	log.Printf("Server started on port %s \n", port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%s", port), nil))
}
