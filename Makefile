include .env
export

BIN_DIR := $(CURDIR)/bin
export PATH := $(BIN_DIR):$(PATH)

GOOSE_VERSION := v3.27.1
GOLANGCI_LINT_VERSION := v2.12.2
MOCKERY_VERSION := v3.5.0

VERSION := $(shell git describe --tags --always)
COMMIT := $(shell git rev-parse --short HEAD)

.PHONY: deps lint-deps lint migrate-up clean build up down debug debug-down
# можно проверять версии перед установкой
deps:
	GOBIN=$(BIN_DIR) go install github.com/pressly/goose/v3/cmd/goose@$(GOOSE_VERSION)

lint-deps:
	# установить golangci-lint
	GOBIN=$(BIN_DIR) go install github.com/golangci/golangci-lint/v2/cmd/golangci-lint@$(GOLANGCI_LINT_VERSION)
	# установить mockery
	GOBIN=$(BIN_DIR) go install github.com/vektra/mockery/v2@$(MOCKERY_VERSION)

lint: lint-deps
	golangci-lint run

test:
	go test ./... -v

test-cover:
	go test ./... -coverprofile=coverage.out
	go tool cover -func=coverage.out
	go tool cover -html=coverage.out

test-integration:
	goose -dir $(MIGRATIONS_DIR) postgres "postgres://$(DB_USER):$(DB_PASSWORD)@localhost:$(DB_PORT)/$(TEST_DB_NAME)?sslmode=$(DB_SSLMODE)" up
	go test -tags=integration ./... -coverprofile=coverage.out -v

migrate-up: deps
	goose -dir $(MIGRATIONS_DIR) postgres "postgres://$(DB_USER):$(DB_PASSWORD)@localhost:$(DB_PORT)/$(DB_NAME)?sslmode=$(DB_SSLMODE)" up

clean:
	rm -rf $(BIN_DIR)

build:
	docker compose build \
		--build-arg APP_VERSION=$(VERSION) \
		--build-arg APP_COMMIT_HASH=$(COMMIT)

up: build
	docker compose up

down:
	docker compose down

debug:
	docker compose \
		-f docker-compose.yml \
		-f docker-compose.debug.yml \
		build \
		--build-arg APP_VERSION=$(VERSION) \
		--build-arg APP_COMMIT_HASH=$(COMMIT)

	docker compose \
		-f docker-compose.yml \
		-f docker-compose.debug.yml \
		up

mock:
	mockery