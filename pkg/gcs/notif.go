package gcs

/*

Sample Payload

{
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
}

*/

// GCSNotification represents the GCS pushed payload
// Capturing only the bits that we need here
type GCSNotification struct {

	// notification
	Kind       string `json:"kind"`
	ID         string `json:"id"`
	SelfLink   string `json:"selfLink"`
	BucketName string `json:"bucket"`

	// object
	ObjectName        string `json:"name"`
	ObjectContentType string `json:"contentType"`
	ObjectSize        string `json:"size"`
	ObjectMD5Hash     string `json:"md5Hash"`
}
