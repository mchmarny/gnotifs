# GCS Notifications

To have GCS push notifications to `gnotifs` on Knative you will need to configure two things:

* [Knative service](../../cmd/service)
* GCS object change notification (this document)

> Assuming you have already set up your `PROJECT_ID`, `KNATIVE_DOMAIN`, and `GCS_KNOWN_PUBLISHER_TOKEN` variable

To start with, define the following environment variables

```shell
# GCS bucket name - can be existing one of one created just for this demo
export GCS_BUCKET_NAME="knative-gcs-sample"
```

### GCS object change notification

There is one aspect of configuring GCS notifications that can't be done from through API. It is the act of adding the Knative domain to the allowed domain list in GCP project. The complete list of steps is outlined [here](https://cloud.google.com/storage/docs/object-change-notification#_Authorize_Endpoint), but the basic steps are:

#### Endpoint Authorization

1. Navigate to [domain verification tab](https://console.cloud.google.com/apis/credentials/domainverification) on the Credentials page in GCP Console (make sure you are in the same GCP project as the GCS bucket from which you want to receive notifications)
3. Click Add domain
4. Enter the service domain (`echo "notif.demo.${KNATIVE_DOMAIN}"`)
5. And click the Add domain button to confirm

If you experience any issues there are few troubleshooting tips [here](https://cloud.google.com/storage/docs/object-change-notification#_Authorize_Endpoint)

#### Service Account

If you don't have a service account already, or if want to create one specific for the GCS notifications, follow the instructions [here](https://cloud.google.com/storage/docs/object-change-notification#_Service_Account). Once you have the service account key, authenticate `gcloud` with that service account.

```shell
gcloud auth activate-service-account $DEMO_SA_EMAIL --key-file $DEMO_SA_KEY_PATH
```
The response from the above command should look like this

```shell
Activated service account credentials for: [YOUR_SA_USERNAME@YOUR_PROJECT.iam.gserviceaccount.com]
```

#### Create notification

The final step is to create a notification

```shell
gsutil notification watchbucket -t $GCS_KNOWN_PUBLISHER_TOKEN -i gnotifs-gcs \
    "https://notif.demo.${KNATIVE_DOMAIN}/gcs" "gs://${GCS_BUCKET_NAME}"
```

The response from the above command should look like this

```shell
Watching bucket gs://$GCS_BUCKET_NAME/ with application URL https://notif.demo.KNATIVE_DOMAIN/gcs ...
Successfully created watch notification channel.
Watch channel identifier: gnotifs-gcs
Canonicalized resource identifier: 35OA-OShWsXoAIxNN8YxxO_7bzw
Client state token: 86b361a2c3c1234f
```

To verify that the notification was configured correctly you can use the list command

```shell
gsutil notification list -o gs://$GCS_BUCKET_NAME
```

The response will look something like this

```shell
...
Notification channel 1:
    Channel identifier: gnotifs-gcs
    Resource identifier: 35OA-OShWsXoAIxNN8YppO_1cft
    ...
```

#### Stop notifications

To stop notifications run this command

```shell
gsutil notification stopchannel gnotifs-gcs ["Resource identifier from the list command"]
```

### Demo

Upload file to your bucket in either browser or using `gsutil` and [see the Knative log output](https://github.com/mchmarny/gnotifs/#demo) the notification data

```shell
echo "Knative test" > ./test.txt
gsutil cp ./test.txt gs://$GCS_BUCKET_NAME
rm ./test.txt
```

## Disclaimer

This is my personal project and it does not represent my employer. I take no responsibility for issues caused by this code. I do my best to ensure that everything works, but if something goes wrong, my apologies is all you will get.


