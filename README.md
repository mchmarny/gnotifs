# kgcs


GCS notifications processing on Knative


## Setup

### Credentials

```shell
gcloud auth activate-service-account \
    mchmarny-laptop-sa@s9-demo.iam.gserviceaccount.com \
    --key-file /Users/mchmarny/.gcp-keys/s9-demo.json
```

Make secure token

```shell
export KGCS_KNOWN_PUBLISHER_TOKEN="$(openssl rand -base64 8 |md5 |head -c16;echo)"
```

Create a bucket

```shell
# TODO: make bucket
export GCS_BUCKET_NAME="knative-gcs-sample"
```

Create notification

```shell
gsutil notification watchbucket \
    -t $KGCS_KNOWN_PUBLISHER_TOKEN \
    -i kgcs-v1 \
    https://kgcs.default.knative.tech/gcs gs://$GCS_BUCKET_NAME
```

List notifications

```shell
gsutil notification list -o gs://$GCS_BUCKET_NAME
```

The response will look something like this

```shell
...
    Notification channel 1:
		Channel identifier: kgcs-v1
		Resource identifier: 35OA-OShWsXoAIxNN8YppO_1cft
        ...
```

Stop notifications

```shell
gsutil notification stopchannel kgcs-v1 ["Resource identifier from list"]
```


## Deploy
