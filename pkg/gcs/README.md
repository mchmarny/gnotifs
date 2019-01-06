# GCS

To configure this demo will will need to configure two things:

* [Knative service](../../cmd/service)
* GCS object change notification (below)

To make things easier, define the following environment variables

```shell
# GCS bucket name
# Can be existing one of one created just for this demo
export GCS_BUCKET_NAME="gnotifs-demo"
```

### GCS object change notification

There is one aspect of configuring GCS notifications that can't be done from through API, it is the act of adding the domain configured on Knative cluster to the allow domain list in GCP project. The complete list of steps is outlined [here](https://cloud.google.com/storage/docs/object-change-notification#_Authorize_Endpoint) but it boils down to following steps.

> Note, if your domain is NOT managed by Google you will also have to [verify the domain ownership](https://cloud.google.com/endpoints/docs/openapi/verify-domain-name)

#### Endpoint Authorization

1. Go to [domain verification tab](https://console.cloud.google.com/apis/credentials/domainverification?_ga=2.186591593.-1146811178.1546727070) on the Credentials page in GCP Console
2. Make sure you are in the same GCP project as the GCS bucket from which you want to receive notifications
3. Click Add domain
4. Enter the service domain (`gnotif.default.${KNATIVE_DOMAIN}`)
5. And click the Add domain button to confirm

If you experience any issues there are few troubleshooting tips [here](https://cloud.google.com/storage/docs/object-change-notification#_Authorize_Endpoint)

#### Service Account

If you don't have a service account already or want to create one specific for the GCS notifications follow the instructions [here](https://cloud.google.com/storage/docs/object-change-notification#_Service_Account). Once you have the service account key, authenticate gcloud with that service account

```shell
gcloud auth activate-service-account \
    YOUR_SERVICE_ACCOUNT_NAME@sPROJECT_ID.iam.gserviceaccount.com \
    --key-file PATH_TO_YOUR_SERVICE_ACCOUNT_KEY.json
```
The response from the above command should look like this

```shell
Activated service account credentials for: [YOUR_SERVICE_ACCOUNT_NAME@sPROJECT_ID.iam.gserviceaccount.com]
```

#### Create notification

The final step is to create a notification

```shell
gsutil notification watchbucket -t $GCS_KNOWN_PUBLISHER_TOKEN -i gnotifs-gcs \
    https://gnotifs.default.$KNATIVE_DOMAIN/gcs gs://$GCS_BUCKET_NAME
```

The response from the above command should look like this

```shell
Watching bucket gs://$GCS_BUCKET_NAME/ with application URL https://gnotif.default.KNATIVE_DOMAIN/gcs ...
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

Upload file to the `$GCS_BUCKET_NAME` in either browser or using `gsutil` see the Knative log output the notification data

## Disclaimer

This is my personal project and it does not represent my employer. I take no responsibility for issues caused by this code. I do my best to ensure that everything works, but if something goes wrong, my apologies is all you will get.


