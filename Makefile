.PHONY: run build test

run:
	go run ./cmd/server

build:
	mkdir -p bin
	go build -o bin/franfit ./cmd/server

test:
	go test ./...
