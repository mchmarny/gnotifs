package cal

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

const (
	testNotificationContent = `{ "kind": "drive#changes" }`
)

func TestNotificationHandler(t *testing.T) {

	req, _ := http.NewRequest("POST", "/drive", strings.NewReader(testNotificationContent))

	req.Header.Add("X-Goog-Channel-Token", knownPublisherToken)
	req.Header.Add("X-Goog-Channel-Id", "gnotifs")
	req.Header.Add("X-Goog-Channel-Expiration", "Tue, 19 Nov 2020 01:13:52 GMT")
	req.Header.Add("X-Goog-Resource-Id", "test")
	req.Header.Add("X-Goog-Resource-State", "test")
	req.Header.Add("X-Goog-Resource-Uri", "test")
	req.Header.Add("X-Goog-Message-Number", "1")

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(NotificationHandler)
	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusAccepted {
		t.Errorf("Handler returned wrong status code: got %v expected %v",
			status, http.StatusAccepted)
		return
	}

}
