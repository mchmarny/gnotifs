package main

import (
	"fmt"
	"log"
	"net/http"
	"encoding/json"

	"github.com/mchmarny/kgcs/pkg/gcs"
	"github.com/mchmarny/kgcs/pkg/utils"
)

func main() {

	log.Print("Starting server...")

	// handlers
	http.HandleFunc("/", defaultHandler)
	http.HandleFunc("/gcs", gcs.NotificationHandler)
	http.HandleFunc("/_healthz", func(w http.ResponseWriter, _ *http.Request) {
		fmt.Fprint(w, "ok")
	})

	// variables
	port := utils.MustGetEnv("PORT", "8080")

	log.Printf("Server started on port %s \n", port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%s", port), nil))
}

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