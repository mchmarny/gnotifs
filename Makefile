GCS_SERVICE_ENDPOINT=https://gnotifs.default.knative.tech/gcs
DRIVE_SERVICE_ENDPOINT=https://gnotifs.default.knative.tech/drive

# DEV

all: test

test:
	go test ./... -v

deps:
	go mod tidy


# DEPLOY

image:
	gcloud builds submit \
		--project=$(PROJECT_ID) \
		--tag gcr.io/$(PROJECT_ID)/gnotifs:latest .

secret:
	kubectl create secret generic gnotifs \
		--from-literal=DRIVE_KNOWN_PUBLISHER_TOKEN=$(DRIVE_KNOWN_PUBLISHER_TOKEN) \
		--from-literal=GCS_KNOWN_PUBLISHER_TOKEN=$(GCS_KNOWN_PUBLISHER_TOKEN)

secret-delete:
	kubectl delete secret gnotifs

deploy:
	kubectl apply -f deployments/gnotifs.yaml

cleanup:
	kubectl delete -f deployments/gnotifs.yaml
	kubectl delete secret gnotif

# GCS

gcs-notif:
	gsutil notification watchbucket \
		-t $(KNOWN_PUBLISHER_TOKEN) \
		-i gnotifs \
		$(GCS_KNOWN_PUBLISHER_TOKEN) gs://$(GCS_BUCKET_NAME)


# DIRVE
