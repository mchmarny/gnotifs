package gcs

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/mchmarny/gnotifs/pkg/utils"
)

var (
	knownPublisherToken = utils.MustGetEnv("KGCS_KNOWN_PUBLISHER_TOKEN", "")
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

// NotificationHandler handles GCS submissions
// https://cloud.google.com/storage/docs/gsutil/commands/notification
func NotificationHandler(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-Type", "application/json")

	// check method
	if r.Method != http.MethodPost {
		log.Printf("wring method: %s", r.Method)
		http.Error(w, "Invalid method. Only POST supported", http.StatusMethodNotAllowed)
		return
	}

	// required
	t := utils.GetHeader("X-Goog-Channel-Token", r)
	s := utils.GetHeader("X-Goog-Resource-State", r)

	// print only others
	utils.GetHeader("X-Goog-Channel-Id", r)
	utils.GetHeader("X-Goog-Resource-Id", r)
	utils.GetHeader("X-Goog-Resource-Uri", r)

	// check for presense/validity of publisher token
	if t != knownPublisherToken {
		log.Printf("Invalid token: %s", t)
		// BadRequest == no retries from GCS
		http.Error(w, fmt.Sprintf("Invalid request (token: %s)", t),
			http.StatusBadRequest)
		return
	}

	// if status is sync then there is no body
	if s == "sync" {
		log.Println("Sync message")
		w.WriteHeader(http.StatusNoContent)
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

	// print JSON for debugging
	// log.Println(string(pb))

	// parse payload
	n := ObjectNotification{}
	if err := json.Unmarshal(pb, &n); err != nil {
		log.Printf("Error decoding notification: %v", err)
		// could be our parsing issue here so BadGateway, GCS will retry
		http.Error(w, fmt.Sprintf("Error decoding notification: %s", err), http.StatusBadGateway)
		return
	}

	// TODO: do something usefull with the pushed message here
	log.Printf("Payload: %v", n)

	// response with accepted status
	w.WriteHeader(http.StatusAccepted)

}
