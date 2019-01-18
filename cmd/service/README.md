# Knative Service

To configure `gnotifs` on Knative you will need to:

* Build image
* Create secret
* Deploy service

> Note, GCP notifications can be only sent to HTTPS endpoints. If your Knative cluster has not been configured with TLS yet, follow the instructions [here](https://github.com/knative/docs/blob/master/serving/using-an-ssl-cert.md) first.

To make things easier, define the following environment variables

```shell
# GCP project ID
export PROJECT_ID="Your GCP project ID"

# Knative domain (only the root e.g. demo.com)
export KNATIVE_DOMAIN="Knative domain here"

# Create tokens which will be shared between GCP and Knative service
# to ensure only valid notifications are processed
export GCS_KNOWN_PUBLISHER_TOKEN="$(openssl rand -base64 16 | md5)"
export DRIVE_KNOWN_PUBLISHER_TOKEN="$(openssl rand -base64 16 | md5)"
```

#### Build image

To build `gnotifs` image using GCP Cloud Build run this command from the root of the project

```shell
gcloud builds submit --project=$(PROJECT_ID) \
    --tag gcr.io/$(PROJECT_ID)/gnotif:latest .
```

The Cloud Build service is pretty verbose, eventually though you should see something like this

```shell
ID           CREATE_TIME          DURATION  SOURCE                                   IMAGES                      STATUS
6905dd3a...  2019-01-23T03:48...  1M43S     gs://PROJECT_cloudbuild/source/15...tgz  gcr.io/PROJECT/gnotifs      SUCCESS
```

Copy the image URI from `IMAGE` column (e.g. `gcr.io/PROJECT_ID/gnotifs`).

#### Create secret

```shell
kubectl create secret generic gnotifs -n demo \
	--from-literal=DRIVE_KNOWN_PUBLISHER_TOKEN=$DRIVE_KNOWN_PUBLISHER_TOKEN \
	--from-literal=GCS_KNOWN_PUBLISHER_TOKEN=$GCS_KNOWN_PUBLISHER_TOKEN
```

The response should be

```shell
secret "gnotifs" created
```

#### Deploy service

Before we can deploy that service to Knative you will need to update the `deployments/gnotifs.yaml` file with the URL of the image we just built.

```yaml
spec:
    container:
        image: gcr.io/PROJECT_ID/gnotifs:latest
```

Once updated, you are ready to publish the `gnotifs` service to Knative cluster

```shell
kubectl apply -f deployments/gnotifs.yaml
```

There response from this command should look someting like this

```shell
service.serving.knative.dev "gnotifs" created
```

You can now test your deployment by invoking its root URL

```shell
open "https://gnotifs.default.${KNATIVE_DOMAIN}"
```

The response should will include the type of notification handlers currently enabled on your Knative service

```json
{
    "handlers": [
        "POST: /gcs",
        "POST: /drive",
    ]
}
```

The `gnotifs` service is now configured on Knative. Go ahead and configure at least one of the [notification sources listed here](https://github.com/mchmarny/gnotifs#gcp-notifications).


## Disclaimer

This is my personal project and it does not represent my employer. I take no responsibility for issues caused by this code. I do my best to ensure that everything works, but if something goes wrong, my apologies is all you will get.


