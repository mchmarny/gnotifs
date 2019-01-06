SERVICE_ENDPOINT=https://kgcs.default.knative.tech/gcs

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
		--tag gcr.io/$(PROJECT_ID)/gnotif:latest .

secret:
	kubectl create secret generic gnotif \
		--from-literal=KNOWN_PUBLISHER_TOKEN=$(KNOWN_PUBLISHER_TOKEN)

delete-secret:
	kubectl delete secret gnotif

deploy:
	kubectl apply -f deployments/gnotif.yaml

cleanup:
	kubectl delete -f deployments/gnotif.yaml
	kubectl delete secret gnotif

# GCS

gcs-notif:
	gsutil notification watchbucket \
		-t $(KNOWN_PUBLISHER_TOKEN) \
		-i gnotif \
		$(SERVICE_ENDPOINT) gs://$(GCS_BUCKET_NAME)


# DIRVE

build-drive-cli:
	export GO111MODULE=on
	go mod tidy
	cd cmd/client/
	go build -o bin/drive-cli

run-drive-cli:
	./bin/drive-cli