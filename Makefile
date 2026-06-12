.PHONY: run mcp test lint migrate seed

run:
	go run ./cmd/api

mcp:
	go run ./cmd/mcp

test:
	go test ./...

lint:
	go vet ./...
	gofmt -d .

migrate:
	go run ./cmd/migrate

seed:
	go run ./cmd/seed --data ./data/seed
