package main

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

const (
	testNotificationContent = `{
		"kind": "storage#object",
		"id": "knative-gcs-sample/main.go/1546720572085362",
		"selfLink": ".../storage/v1/b/BUCKET/o/main.go",
		"name": "main.go",
		"bucket": "knative-gcs-sample",
		"generation": "1546720572085362",
		"metageneration": "1",
		"contentType": "application/octet-stream",
		"timeCreated": "2019-01-05T20:36:12.085Z",
		"updated": "2019-01-05T20:36:12.085Z",
		"timeStorageClassUpdated": "2019-01-05T20:36:12.085Z",
		"size": "1787",
		"md5Hash": "Sh8fXcL5oD8o6va5zE5BWg==",
		"mediaLink": "...download/storage/v1/b/BUCKET/o/main.go?generation=154&alt=media",
		"crc32c": "KyB3Ww==",
		"etag": "CPKokJK/198CEAE="
	}`
)

func TestNotificationHandlerWithValidToken(t *testing.T) {

	testToken := mustGetEnv("KGCS_KNOWN_PUBLISHER_TOKEN", "")

	req, _ := http.NewRequest("POST", "/gcs", strings.NewReader(testNotificationContent))

	req.Header.Add("X-Goog-Channel-Token", testToken)
	req.Header.Add("X-Goog-Channel-Id", "kgcs")
	req.Header.Add("X-Goog-Resource-Id", "test")
	req.Header.Add("X-Goog-Resource-State", "test")
	req.Header.Add("X-Goog-Resource-Uri", "test")

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(notificationHandler)
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
	handler := http.HandlerFunc(notificationHandler)
	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusBadRequest {
		t.Errorf("Handler returned invalid status for bad token: got %v expected %v",
			status, http.StatusBadRequest)
		return
	}

}
