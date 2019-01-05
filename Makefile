SERVICE_ENDPOINT=https://kgcs.default.knative.tech/gcs

all: test

test:
	go test ./... -v

run:
	go run main.go utils.go handlers.go notif.go

deps:
	go mod tidy

notif:
	gsutil notification watchbucket \
		-t $(KGCS_KNOWN_PUBLISHER_TOKEN) \
		-i kgcs \
		$(SERVICE_ENDPOINT) gs://$(GCS_BUCKET_NAME)

image:
	gcloud builds submit \
		--project=$(PROJECT_ID) \
		--tag gcr.io/$(PROJECT_ID)/kgcs:latest .

secret:
	kubectl create secret generic kgcs \
		--from-literal=KGCS_KNOWN_PUBLISHER_TOKEN=$(KGCS_KNOWN_PUBLISHER_TOKEN)

delete-secret:
	kubectl delete secret kgcs

deploy:
	kubectl apply -f service.yaml

cleanup:
	kubectl delete -f service.yaml
	kubectl delete secret kgcs