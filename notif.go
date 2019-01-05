package main

import (
	"time"
)

/*

Sample Payload

{
 "kind": "storage#object",
 "id": "BucketName/ObjectName",
 "selfLink": "https://www.googleapis.com/storage/v1/b/BucketName/o/ObjectName",
 "name": "ObjectName",
 "bucket": "BucketName",
 "generation": "1367014943964000",
 "metageneration": "1",
 "contentType": "application/octet-stream",
 "updated": "2013-04-26T22:22:23.832Z",
 "size": "10",
 "md5Hash": "xHZY0QLVuYng2gnOQD90Yw==",
 "mediaLink": "https://www.googleapis.com/storage/v1/b/BucketName/o/ObjectName?generation=1367014943964000&alt=media",
 "owner": { "entity": "user-jane@gmail.com" },
 "crc32c": "C7+82w==",
 "etag": "COD2jMGv6bYCEAE="
}

*/

// Notification represents the GCS pushed payload
type Notification struct {
	ID          string             `json:"id"`
	Kind        string             `json:"kind"`
	SelfLink    string             `json:"selfLink"`
	ObjectName  string             `json:"name"`
	BucketName  string             `json:"bucket"`
	ContentType string             `json:"contentType"`
	Updated     time.Time          `json:"updated"`
	Size        string             `json:"size"`
	MD5Hash     string             `json:"md5Hash"`
	Owner       *NotificationOwner `json:"owner"`
}

// NotificationOwner represents the owner object of Notification
type NotificationOwner struct {
	Entity string `json:"entity"`
}
