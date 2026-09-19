GO ?= go
GOFILES := $(shell find . -type f -name '*.go' -not -path './internal/httpapi/generated/*' -not -path './internal/postgres/sqlc/*')

.PHONY: dev generate generate-check build test lint check migrate

dev:
	$(GO) tool air

generate:
	$(GO) tool oapi-codegen --config oapi-codegen.yaml ../../packages/contracts/openapi.yaml
	$(GO) tool sqlc generate

generate-check: generate
	git diff --exit-code -- internal/httpapi/generated internal/postgres/sqlc

build:
	$(GO) build ./cmd/api ./cmd/admin ./cmd/migrate

test:
	$(GO) tool gotestsum -- ./...

lint:
	test -z "$$($(GO)fmt -l $(GOFILES))"
	$(GO) vet ./...

check: generate-check lint test build

migrate:
	$(GO) run ./cmd/migrate
