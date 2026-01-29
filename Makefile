APP_NAME := test-messenger
APP_PORT := 8080

.PHONY: help
help:
	@echo "Available commands:"
	@echo "  make run        - run app locally (uses env vars)"
	@echo "  make test       - run tests"
	@echo "  make build      - build Go binary"
	@echo "  make lint       - run go vet"
	@echo "  make up         - start app with docker-compose"
	@echo "  make down       - stop docker-compose"
	@echo "  make migrate    - run goose migrations"
	@echo "  make migrate-down - rollback goose migrations"

.PHONY: run
run:
	go run ./cmd/app

.PHONY: test
test:
	go test ./...

.PHONY: build
build:
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o bin/$(APP_NAME) ./cmd/app

.PHONY: lint
lint:
	go vet ./...

.PHONY: up
up:
	docker compose up --build

.PHONY: down
down:
	docker compose down -v

.PHONY: migrate
migrate:
	goose -dir migrations postgres "$$DATABASE_DSN" up

.PHONY: migrate-down
migrate-down:
	goose -dir migrations postgres "$$DATABASE_DSN" down
