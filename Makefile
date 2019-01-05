GCP_PROJECT_NAME=s9-demo
BINARY_NAME=kgcs

all: test

run:
	go run main.go

deps:
	go mod tidy

image:
	gcloud builds submit \
		--project=$(GCP_PROJECT_NAME) \
		--tag gcr.io/$(GCP_PROJECT_NAME)/$(BINARY_NAME):latest .

secret:
	kubectl create secret generic $(BINARY_NAME) \
		--from-literal=KNOWN_PUBLISHER_TOKEN=$(KNOWN_PUBLISHER_TOKEN)

deploy:
	kubectl apply -f service.yaml