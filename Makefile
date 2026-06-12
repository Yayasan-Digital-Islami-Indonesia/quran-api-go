.PHONY: run mcp test lint migrate seed swag

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

swag:
	swag init -g cmd/api/main.go -o docs --outputTypes go,yaml
	mkdir -p docs/api-reference
	cp docs/swagger.yaml docs/api-reference/openapi.yaml
