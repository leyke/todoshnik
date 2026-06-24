include .env
export

BIN_DIR := $(CURDIR)/bin
export PATH := $(BIN_DIR):$(PATH)

GOOSE_VERSION := v3.27.1
GOLANGCI_LINT_VERSION := v2.12.2

.PHONY: deps lint migrate-up

deps:
	GOBIN=$(BIN_DIR) go install github.com/pressly/goose/v3/cmd/goose@$(GOOSE_VERSION)
	GOBIN=$(BIN_DIR) go install github.com/golangci/golangci-lint/v2/cmd/golangci-lint@$(GOLANGCI_LINT_VERSION)

lint: deps
	golangci-lint run

migrate-up: deps
	goose -dir $(MIGRATIONS_DIR) postgres "postgres://$(DB_USER):$(DB_PASSWORD)@localhost:$(DB_PORT)/$(DB_NAME)?sslmode=$(DB_SSLMODE)" up

clean:
	rm -rf $(BIN_DIR)

build:
	docker compose up --build