# kgcs

GCS notifications processing using Knative service demo

GCP exposes object changed notification. On GCP you can easy wire it to PubSub of even your own code on GCF. If you want to wire these notification to external endpoint like service deployed on Knative there are few additional steps. This demo will walk through the entire process.

> Note, this project will be refactored to support other notification types like Google Drive

## Setup

To configure this demo will will need to configure two things:

* [Knative service]()
* [GCS object change notification]()

To make things easier, define the following enviernment variables

```shell
# GCP project ID
export PROJECT_ID="Your project ID"

# Knative domain (only the root e.g. demo.com)
export KNATIVE_DOMAIN="demo.com"

# Token which will be shared between GCS and Knative service
# to ensure only GCS notifications are processed
export KGCS_KNOWN_PUBLISHER_TOKEN="$(openssl rand -base64 16 |md5 |head -c16;echo)"

# GCS bucket name
# Can be existing one of one created just for this demo
export KGCS_BUCKET_NAME="knative-gcs-demo"
```

### Knative service

To deploy this service to Knative we need to:

* Build image
* Create secret
* Deploy service

> Note, GCS notifications can be only sent to HTTPS endpoints. If your Knative cluster has not been configured with TLS yet, follow the instructions [here](https://github.com/knative/docs/blob/master/serving/using-an-ssl-cert.md) first.

#### Build image

From the root of this project run

```shell
gcloud builds submit \
    --project=$(PROJECT_ID) \
    --tag gcr.io/$(PROJECT_ID)/kgcs:latest .
```

The build service is pretty verbose in output, eventually though you should see something like this

```shell
ID           CREATE_TIME          DURATION  SOURCE                                   IMAGES                      STATUS
6905dd3a...  2019-01-23T03:48...  1M43S     gs://PROJECT_cloudbuild/source/15...tgz  gcr.io/PROJECT/kgcs SUCCESS
```

Copy the image URI from `IMAGE` column (e.g. `gcr.io/PROJECT_ID/kgcs`).

#### Create secret

```shell
kubectl create secret generic kgcs \
	--from-literal=KGCS_KNOWN_PUBLISHER_TOKEN=$(KGCS_KNOWN_PUBLISHER_TOKEN)
```

The response should be

```shell
secret "kgcs" created
```

#### Deploy service

Before we can deploy that service to Knative you will need to update the `service.yaml` file to the image URL you captured from the `Build image` step.

```yaml
spec:
    container:
        image: gcr.io/PROJECT_ID/kgcs:latest
```

To test your deployment you should be able to invoke the root of the `kgcs` and see

```json
{
    "handlers": [ "POST: /gcs" ]
}
```

That means the Knative service found the necessary secret and it ready to process notifications from GCS

### GCS object change notification

There is one aspect of configuring GCS notifications that can't be done from through API, it is the act of adding the domain configured on Knative cluster to the allow domain list in GCP project. The complete list of steps is outlined [here](https://cloud.google.com/storage/docs/object-change-notification#_Authorize_Endpoint) but it boils down to following steps.

> Note, if your domain is NOT managed by Google you will also have to [verify the domain ownership](https://cloud.google.com/endpoints/docs/openapi/verify-domain-name)

#### Endpoint Authorization

1. Go to [domain verification tab](https://console.cloud.google.com/apis/credentials/domainverification?_ga=2.186591593.-1146811178.1546727070) on the Credentials page in GCP Console
2. Make sure you are in the same GCP project as the GCS bucket from which you want to receive notifications
3. Click Add domain
4. Enter the service domain (`kgcs.default.${KNATIVE_DOMAIN}`)
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
gsutil notification watchbucket -t $KGCS_KNOWN_PUBLISHER_TOKEN -i kgcs \
    https://kgcs.default.$KNATIVE_DOMAIN/gcs gs://$KGCS_BUCKET_NAME
```

The response from the above command should look like this

```shell
Watching bucket gs://$KGCS_BUCKET_NAME/ with application URL https://kgcs.default.KNATIVE_DOMAIN/gcs ...
Successfully created watch notification channel.
Watch channel identifier: kgcs
Canonicalized resource identifier: 35OA-OShWsXoAIxNN8YxxO_7bzw
Client state token: 86b361a2c3c1234f
```

To verify that the notification was configured correctly you can use the list command

```shell
gsutil notification list -o gs://$KGCS_BUCKET_NAME
```

The response will look something like this

```shell
...
Notification channel 1:
    Channel identifier: kgcs
    Resource identifier: 35OA-OShWsXoAIxNN8YppO_1cft
    ...
```

#### Stop notifications

To stop notifications run this command

```shell
gsutil notification stopchannel kgcs ["Resource identifier from the list command"]
```

### Demo

The `kgcs` app doesn't do much with the submitted GCS notifications other than log them. To demo this app you can query the Kubernetes logs.

```shell
kubectl -l 'serving.knative.dev/service=kgcs' logs -c user-container
```

Then in the browser you can either upload or delete a file in the `$KGCS_BUCKET_NAME` bucket and see the console output the notification data.


## Disclaimer

This is my personal project and it does not represent my employer. I take no responsibility for issues caused by this code. I do my best to ensure that everything works, but if something goes wrong, my apologies is all you will get.


