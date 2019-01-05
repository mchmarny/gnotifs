package main

import (
	"fmt"
	"log"
	"net/http"
)

var (
	knownPublisherToken = ""
)

func main() {

	log.Print("Starting server...")

	// handlers
	http.HandleFunc("/", defaultHandler)
	http.HandleFunc("/gcs", notificationHandler)
	http.HandleFunc("/_healthz", func(w http.ResponseWriter, _ *http.Request) {
		fmt.Fprint(w, "ok")
	})

	// variables
	knownPublisherToken = mustGetEnv("KNOWN_PUBLISHER_TOKEN", "")
	port := mustGetEnv("PORT", "8080")

	log.Printf("Server started on port %s \n", port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%s", port), nil))
}
