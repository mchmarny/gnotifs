# gnotifs

GCP notifications processing using Knative service demo

GCP exposes many types of evens as notifications. You can easy wire it to PubSub of even your own code on GCF. If you want to wire these notifications to external endpoint like service deployed on Knative there are few additional steps. This demo will walk through the entire process.

## Knative Service

The instructions on building and deploying the `gnotifs` service onto Knative are located [here](cmd/service).


## Notifications

Once the notification target service is configured on Knative, you will have to configure the event sources (triggers) individually. The currently supported notifications are:

* [Google Cloud Storage (GCS)](pkg/gcs)
* [Google Drive (Drive)](pkg/drive)


## Disclaimer

This is my personal project and it does not represent my employer. I take no responsibility for issues caused by this code. I do my best to ensure that everything works, but if something goes wrong, my apologies is all you will get.


