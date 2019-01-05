package main

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

// Notification represents the GCS pushed payload
// Capturing only the bits that we need here
type Notification struct {
	Kind        string `json:"kind"`
	ID          string `json:"id"`
	SelfLink    string `json:"selfLink"`
	ObjectName  string `json:"name"`
	BucketName  string `json:"bucket"`
	ContentType string `json:"contentType"`
	Size        string `json:"size"`
	MD5Hash     string `json:"md5Hash"`
}
