# Knative Service

To configure `gnotifs` on Knative you will need to:

* Build image
* Create secret
* Deploy service

> Note, GCP notifications can be only sent to HTTPS endpoints. If your Knative cluster has not been configured with TLS yet, follow the instructions [here](https://github.com/knative/docs/blob/master/serving/using-an-ssl-cert.md) first.

To make things easier, define the following environment variables

```shell
# GCP project ID
export PROJECT_ID="Your project ID"

# Knative domain (only the root e.g. demo.com)
export KNATIVE_DOMAIN="demo.com"

# Token which will be shared between GCP and Knative service
# to ensure only valid notifications are processed
export GCS_KNOWN_PUBLISHER_TOKEN="$(openssl rand -base64 16 |md5 |head -c16;echo)"
export DRIVE_KNOWN_PUBLISHER_TOKEN="$(openssl rand -base64 16 |md5 |head -c16;echo)"

```


#### Build image

From the root of this project run

```shell
gcloud builds submit \
    --project=$(PROJECT_ID) \
    --tag gcr.io/$(PROJECT_ID)/gnotif:latest .
```

The build service is pretty verbose in output, eventually though you should see something like this

```shell
ID           CREATE_TIME          DURATION  SOURCE                                   IMAGES                      STATUS
6905dd3a...  2019-01-23T03:48...  1M43S     gs://PROJECT_cloudbuild/source/15...tgz  gcr.io/PROJECT/gnotif SUCCESS
```

Copy the image URI from `IMAGE` column (e.g. `gcr.io/PROJECT_ID/gnotif`).

#### Create secret

```shell
kubectl create secret generic gnotif \
	--from-literal=KNOWN_PUBLISHER_TOKEN=$(KNOWN_PUBLISHER_TOKEN)
```

The response should be

```shell
secret "gnotif" created
```

#### Deploy service

Before we can deploy that service to Knative you will need to update the `service.yaml` file to the image URL you captured from the `Build image` step.

```yaml
spec:
    container:
        image: gcr.io/PROJECT_ID/gnotif:latest
```

To test your deployment you should be able to invoke the root of the `gnotif` and see

```json
{
    "handlers": [
        "POST: /gcs",
        "POST: /drive",
    ]
}
```

That means the Knative service found the necessary secret and it ready to process notifications from GCS

Go ahead and configure at least one of the notification sources listed [here](https://github.com/mchmarny/gnotifs)

### Demo

The `gnotif` service doesn't do much with the submitted notifications. It does log them so we can tail on the Kubernetes logs to demo

```shell
kubectl -l 'serving.knative.dev/service=gnotif' logs -c user-container
```


## Disclaimer

This is my personal project and it does not represent my employer. I take no responsibility for issues caused by this code. I do my best to ensure that everything works, but if something goes wrong, my apologies is all you will get.


