# Makefile for the Go backend

.PHONY: all build test lint format clean

# Install golangci-lint with:
# go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
# Or check the official installation guide: https://golangci-lint.run/usage/install/

# Commands
GOLINT = golangci-lint
GOFMT = gofmt
GOTEST = go test
GOBUILD = go build

# Linting and formatting targets
LINT_TARGETS := ./...
FORMAT_TARGETS := ./

# Build artifacts
BINARY_NAME = backend
BUILD_OUTPUT = bin/$(BINARY_NAME)

all: build

build: format lint
	@echo "==> Building application..."
	mkdir -p bin
	$(GOBUILD) -o $(BUILD_OUTPUT) ./cmd/api/main.go

test: format lint
	@echo "==> Running tests..."
	$(GOTEST) $(LINT_TARGETS)

# If golangci-lint is not installed, linting will fail. Please install it to run linting.
# go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
lint:
	@echo "==> Running linter..."

	$(GOLINT) run $(LINT_TARGETS)

format:
	@echo "==> Formatting code..."
	$(GOFMT) -w $(FORMAT_TARGETS)

clean:
	@echo "==> Cleaning up..."
	rm -f $(BUILD_OUTPUT)
