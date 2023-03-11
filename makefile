SHELL := /bin/bash

.PHONY: all build test deps deps-cleancache

GOCMD=go
BUILD_DIR=build
BINARY_DIR=$(BUILD_DIR)/bin
CODE_COVERAGE=code-coverage

all: test build


run: ## Start application
	$(GOCMD) run ./cmd/api

wire: ## Generate wire_gen.go
	cd pkg/di && wire

#swag: ## Generate swagger docs
#	swag init -g pkg/http/handler/user.go -o ./cmd/api/docs

## swag init -g pkg/api/handler/user.go -o ./cmd/api/docs
swag: ## Generate swagger docs
	cd cmd/api && swag init --parseDependency --parseInternal --parseDepth 1 -md ./documentation -o ./docs

