package drive

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/mchmarny/gnotifs/pkg/utils"
)

var (
	knownPublisherToken = utils.MustGetEnv("DRIVE_KNOWN_PUBLISHER_TOKEN", "")
)

/*

POST
Content-Type: application/json; charset="utf-8"
X-Goog-Channel-ID: channel-ID-value
X-Goog-Channel-Token: channel-token-value
X-Goog-Channel-Expiration: expiration-date-and-time // In human-readable format; present only if channel expires.
X-Goog-Resource-ID: identifier-for-the-watched-resource
X-Goog-Resource-URI: version-specific-URI-of-the-watched-resource
X-Goog-Resource-State: sync
X-Goog-Message-Number: 1
*/

// NotificationHandler handles Google Drive notifications
// https://developers.google.com/drive/api/v3/push
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
	// print only others
	utils.GetHeader("X-Goog-Channel-Id", r)
	utils.GetHeader("X-Goog-Resource-Id", r)
	utils.GetHeader("X-Goog-Resource-Uri", r)
	utils.GetHeader("X-Goog-Channel-Expiration", r)
	utils.GetHeader("X-Goog-Changed", r)
	utils.GetHeader("X-Goog-Message-Number", r)

	// check for presense/validity of publisher token
	if t != knownPublisherToken {
		log.Printf("Invalid token: %s", t)
		// BadRequest == no retries from GCS
		http.Error(w, fmt.Sprintf("Invalid request (token: %s)", t),
			http.StatusBadRequest)
		return
	}

	// print document state
	log.Printf("Document state: %s", s)

	// get payload
	pb, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Printf("Error capturing payload: %v", err)
	} else {
		// print JSON for debugging
		log.Printf("BODY: %s", string(pb))
	}

	// response with accepted status
	w.WriteHeader(http.StatusAccepted)

}
