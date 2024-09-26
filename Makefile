.PHONY: format vet build run

format:
	@go fmt ./...

vet:
	@go vet ./...

build: format vet
	@go build -o bin/sourcecollector main.go

run: build
	@./bin/sourcecollector --input=$(input) --output=$(output) --fast=$(fast)