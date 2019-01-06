package gcs

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/mchmarny/gnotifs/pkg/utils"
)

const (
	testNotificationContent = `{
		"kind": "storage#object",
		"id": "knative-gcs-sample/test.txt/1546728971575475",
		"selfLink": "https://www.googleapis.com/storage/v1/b/knative-gcs-sample/o/test.txt",
		"name": "test.txt",
		"bucket": "knative-gcs-sample",
		"generation": "1546728971575475",
		"metageneration": "1",
		"contentType": "text/plain",
		"timeCreated": "2019-01-05T22:56:11.575Z",
		"updated": "2019-01-05T22:56:11.575Z",
		"storageClass": "MULTI_REGIONAL",
		"timeStorageClassUpdated": "2019-01-05T22:56:11.575Z",
		"size": "714",
		"md5Hash": "7zmNlW0o3P9yZxz/3NLClw==",
		"mediaLink": "https://www.googleapis.com/download/storage/v1/b/knative-gcs-sample/o/test.txt?generation=1546728971575475&alt=media",
		"crc32c": "Go70Vg==",
		"etag": "CLPBqLfe198CEAE="
 	}`
)

func TestNotificationHandlerWithValidToken(t *testing.T) {

	testToken := utils.MustGetEnv("KGCS_KNOWN_PUBLISHER_TOKEN", "")

	req, _ := http.NewRequest("POST", "/gcs", strings.NewReader(testNotificationContent))

	req.Header.Add("X-Goog-Channel-Token", testToken)
	req.Header.Add("X-Goog-Channel-Id", "kgcs")
	req.Header.Add("X-Goog-Resource-Id", "test")
	req.Header.Add("X-Goog-Resource-State", "test")
	req.Header.Add("X-Goog-Resource-Uri", "test")

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(NotificationHandler)
	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusAccepted {
		t.Errorf("Handler returned wrong status code: got %v expected %v",
			status, http.StatusAccepted)
		return
	}

}

func TestNotificationHandlerWithInvalidToken(t *testing.T) {

	req, _ := http.NewRequest("POST", "/gcs", strings.NewReader(testNotificationContent))

	req.Header.Add("X-Goog-Channel-Token", "invalidToken")
	req.Header.Add("X-Goog-Channel-Id", "kgcs")
	req.Header.Add("X-Goog-Resource-Id", "test")
	req.Header.Add("X-Goog-Resource-State", "test")
	req.Header.Add("X-Goog-Resource-Uri", "test")

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(NotificationHandler)
	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusBadRequest {
		t.Errorf("Handler returned invalid status for bad token: got %v expected %v",
			status, http.StatusBadRequest)
		return
	}

}
