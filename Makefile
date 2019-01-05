
all: test

test:
	go test ./... -v

run:
	go run main.go utils.go handlers.go notif.go

deps:
	go mod tidy

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