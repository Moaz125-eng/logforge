.PHONY: build test run docker

build:
	go build -o bin/logforge-server ./cmd/server
	go build -o bin/logforge-agent ./cmd/agent
	go build -o bin/logforge-bench ./cmd/bench

test:
	go test ./...

run:
	go run ./cmd/server

docker:
	docker compose up --build
