package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

/*

POST
Content-Type: application/json; charset="utf-8"
X-Goog-Channel-Id: ChannelId
X-Goog-Channel-Token: ClientToken
X-Goog-Resource-Id: ResourceId
X-Goog-Resource-State: ResourceState
X-Goog-Resource-Uri: https://www.googleapis.com/storage/v1/b/BucketName/o?alt=json
*/

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

// notificationHandler handles GCS submissions
// https://cloud.google.com/storage/docs/gsutil/commands/notification
func notificationHandler(w http.ResponseWriter, r *http.Request) {

	// check method
	if r.Method != http.MethodPost {
		log.Printf("wring method: %s", r.Method)
		http.Error(w, "Invalid method. Only POST supported", http.StatusMethodNotAllowed)
		return
	}

	// required
	t := getHeader("X-Goog-Channel-Token", r)

	// print only others
	getHeader("X-Goog-Channel-Id", r)
	getHeader("X-Goog-Resource-Id", r)
	getHeader("X-Goog-Resource-State", r)
	getHeader("X-Goog-Resource-Uri", r)

	// check for presense/validity of publisher token
	if t != knownPublisherToken {
		log.Printf("Invalid token: %s", t)
		http.Error(w, fmt.Sprintf("Invalid request (token: %s)", t),
			http.StatusBadRequest)
		return
	}

	// get payload
	pb, err := ioutil.ReadAll(r.Body)
	defer r.Body.Close()
	if err != nil {
		log.Printf("Error capturing payload: %v", err)
		http.Error(w, fmt.Sprintf("Error capturing payload: %s", err), http.StatusBadRequest)
		return
	}

	log.Printf("Body: %s", string(pb))

	// parse payload
	n := Notification{}
	if err := json.Unmarshal(pb, &n); err != nil {
		log.Printf("Error decoding notification: %v", err)
		http.Error(w, fmt.Sprintf("Error decoding notification: %s", err), http.StatusBadRequest)
		return
	}

	// TODO: Check the MD5 of the payload

	// TODO: do something usefull with the pushed message here
	log.Printf("Payload: %s", n)

	// response with accepted status
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusAccepted)

}
