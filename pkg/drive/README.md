# Drive Notifications

Google Drive supports notifications for one document or all files in your drive. In this example will set up notification for one Google Doc document and process them in Knative. To configure this demo will will need to setup two things:

* [Knative service](../../cmd/service)
* Google Drive document change notification (this document)

## Google Drive document change notification

There is one aspect of configuring Google Drive notifications that can't be done from through API, it is the act of adding the domain configured on Knative cluster to the allow domain list in GCP project. The complete list of steps is outlined [here](https://cloud.google.com/storage/docs/object-change-notification#_Authorize_Endpoint)

## Get OAuth Credentials

> Assuming you have already set up your `PROJECT_ID`, `KNATIVE_DOMAIN`, and `DRIVE_KNOWN_PUBLISHER_TOKEN` variable

To create your OAuth credentials, run the below command to navigate in browser to the GCP console using this command:

```shell
open "https://console.developers.google.com/apis/credentials?project=${PROJECT_ID}&authuser=0"
```

Once there:

* Click `Create Credentails`
   * Select `Create OAuth client ID`
   * Select `Other` for `Application type`
   * Enter `gnotifs` as `Name`
   * Click `Create`
* Click `Download JSON` icon on the far right for the `gnotifs` line item

## Get OAuth Token

From the downloaded credentials JSON copy `client_id` and `client_secret`. To make things easier, export these as environment variables like below.

```shell
export NOTIF_OAUTH_CLIENT_ID=""
export NOTIF_OAUTH_CLIENT_SECRET=""
```

The go in browser to generate request code

```shell
open "https://accounts.google.com/o/oauth2/auth?client_id=${NOTIF_OAUTH_CLIENT_ID}&redirect_uri=urn:ietf:wg:oauth:2.0:oob&scope=https://www.googleapis.com/auth/drive.metadata.readonly&response_type=code"
```

Follow the on-screen prompts and when done, copy the `code`. Again, export it as environment variables.

```shell
export NOTIF_OAUTH_AUTH_CODE=
```

Now it's time to get the OAuth token

```shell
curl -d client_id=$NOTIF_OAUTH_CLIENT_ID \
     -d client_secret=$NOTIF_OAUTH_CLIENT_SECRET \
     -d grant_type=authorization_code \
     -d redirect_uri=urn:ietf:wg:oauth:2.0:oob \
     -d code=$NOTIF_OAUTH_AUTH_CODE \
     https://accounts.google.com/o/oauth2/token
```

Copy the `access_token` from the response and export it as environment variables. We are also are goign to generate UUID to help as track Google Drive notification channel.

```shell
export NOTIF_OAUTH_AUTH_TOKEN=
export DRIVE_CHANNEL_ID=$(uuidgen)
export DRIVE_SERVICE_ENDPOINT="https://notif.demo.${KNATIVE_DOMAIN}/drive"
```

## Drive Notification Channel

First create brand new Google Drive document and copy its ID. For example, if the complete Google Docs URL looks like this

```shell
https://docs.google.com/document/d/1uUR7hxvQMZOTDiWLg84mmt-Zee_f7_kuX366pQGTkXk/edit
```

The document ID will be in between `d/` and `/edit` (e.g. `1uUR7hxvQMZOTDiWLg84mmt-Zee_f7_kuX377pQGTkXk`). Again, export the document ID as environment variables

```shell
export NOTIFICATION_DOC_ID="1uUR7hxvQMZOTDiWLg84mmt-Zee_f7_kuX377pQGTkXk"
```

To create Drive notification channel now run this command

```shell
curl -X POST -H "Content-Type: application/json" \
    -H "Authorization: Bearer ${NOTIF_OAUTH_AUTH_TOKEN}" \
    -d "{ 'id': '${DRIVE_CHANNEL_ID}', 'type': 'web_hook', 'address': '${DRIVE_SERVICE_ENDPOINT}', 'token': '${DRIVE_KNOWN_PUBLISHER_TOKEN}' }" \
    https://www.googleapis.com/drive/v3/files/${NOTIFICATION_DOC_ID}/watch
```

Capture the `resourceId` if you ever need to stop that channel

> Note, by default, Drive API will be sending these notifications only for 1 hr. If you want to extend it, include an `expiration` property string set to a Unix timestamp (in ms) in the above message. The maximum time for Drive notifications is 1 day (86400 seconds)!


## Stop Drive Channel Notification Channel

The channels expire themselves and you will have to recreate it bit if you want to stop it manually, you can run this command


```shell
curl -X POST -H "Content-Type: application/json" \
    -H "Authorization: Bearer ${NOTIF_OAUTH_AUTH_TOKEN}" \
    -d '{ "id": "${DRIVE_CHANNEL_ID}", "resourceId": "YOUR_CHANNEL_RESOURCE_ID" }' \
    https://www.googleapis.com/drive/v3/channels/stop
```

## Cleanup

```shell
unset OAUTH_CLIENT_ID
unset OAUTH_CLIENT_SECRET
unset OAUTH_AUTH_CODE
unset OAUTH_AUTH_TOKEN
unset DRIVE_CHANNEL_ID
unset DRIVE_SERVICE_ENDPOINT
unset NOTIFICATION_DOC_ID
```
