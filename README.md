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
export KNOWN_PUBLISHER_TOKEN="$(openssl rand -base64 8 |md5 |head -c16;echo)"
```

Create a bucket

```shell
# TODO: make bucket
export GCS_BUCKET_NAME="knative-gcs-sample"
```

Create notification

```shell
gsutil notification watchbucket \
    -t $KNOWN_PUBLISHER_TOKEN \
    -i kgcs \
    https://kgcs.default.knative.tech/gcs gs://$GCS_BUCKET_NAME
```




## Deploy
