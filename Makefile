.PHONY: format-backend vet-backend

format:
	@go fmt ./...

vet:
	@go vet ./...

build: format vet
	@go build -o bin/sourcecollector cmd/sourcecollector/main.go

run: build
	@./bin/sourcecollector