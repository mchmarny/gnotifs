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

func getHeader(key string, r *http.Request) string {
	v := r.Header.Get(key)
	log.Printf("Header[%s] %s", key, v)
	return v
}

// gcsObjectHandler handles GCS submissions
// https://cloud.google.com/storage/docs/gsutil/commands/notification
func gcsObjectHandler(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-Type", "application/json")

	// check method
	if r.Method != http.MethodPost {
		log.Printf("wring method: %s", r.Method)
		http.Error(w, "Invalid method. Only POST supported", http.StatusMethodNotAllowed)
		return
	}

	// required
	chToken := getHeader("X-Goog-Channel-Token", r)

	// other
	getHeader("X-Goog-Channel-Id", r)
	getHeader("X-Goog-Resource-Id", r)
	getHeader("X-Goog-Resource-State", r)
	getHeader("X-Goog-Resource-Uri", r)

	// check for presense/validity of publisher token
	if chToken != knownPublisherToken {
		log.Printf("Invalid token: %s", chToken)
		http.Error(w, fmt.Sprintf("Invalid request (token: %s)", chToken),
			http.StatusBadRequest)
		return
	}

	// get payload
	pb, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Printf("Error capturing payload: %v", err)
		http.Error(w, fmt.Sprintf("Error capturing payload: %s", err), http.StatusBadRequest)
		return
	}

	log.Printf("Body: %s", pb)

	// parse payload
	notif := &Notification{}
	if err := json.NewDecoder(r.Body).Decode(notif); err != nil {
		log.Printf("Error decoding notification: %v", err)
		http.Error(w, fmt.Sprintf("Error decoding notification: %s", err), http.StatusBadRequest)
		return
	}

	// TODO: Check the MD5 of the payload

	// TODO: do something usefull with the pushed message here
	log.Printf("Payload: %s", notif)

	// response with the parsed payload data
	w.WriteHeader(http.StatusAccepted)

}
