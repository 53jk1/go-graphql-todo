# Variables
APP_NAME = go-graphql-todo
VERSION = v1.0.0
GOFLAGS ?= $(GOFLAGS:)
BUILD_DIR ?= ./build
GO_FILES := $(shell find . -type f -name '*.go' -not -path "./vendor/*")

.PHONY: build clean run format test lint gqlgen

build:
	@echo "Building $(APP_NAME) $(VERSION)"
	mkdir -p $(BUILD_DIR)
	go build -ldflags="-X 'main.version=$(VERSION)'" $(GOFLAGS) -o $(BUILD_DIR)/$(APP_NAME) ./cmd/server

clean:
	@echo "Cleaning up..."
	rm -rf $(BUILD_DIR)

run:
	@echo "Running $(APP_NAME) $(VERSION)"
	$(BUILD_DIR)/$(APP_NAME)

run-dev:
	@echo "Running $(APP_NAME) $(VERSION) in development mode"
	go run $(GOFLAGS) ./cmd/server

format:
	@echo "Formatting source code..."
	go fmt $(GO_FILES)

test:
	@echo "Running tests..."
	go test $(GOFLAGS) ./...

lint:
	@echo "Running linter..."
	golangci-lint run

install-gqlgen:
	@echo "Installing gqlgen..."
	go get -u github.com/99designs/gqlgen

gqlgen:
	@echo "Generating GraphQL code..."
	go run github.com/99designs/gqlgen generate

.PHONY: help
help:
	@echo "Available targets:"
	@echo "  build       Build the server application"
	@echo "  clean       Remove build artifacts"
	@echo "  run         Run the server application"
	@echo "  format      Format the source code"
	@echo "  test        Run tests"
	@echo "  lint        Run linter"
	@echo "  gqlgen      Generate GraphQL code"
	@echo "  help        Show this help message"