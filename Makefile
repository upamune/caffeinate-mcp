.PHONY: build test clean install release-snapshot help

# Variables
BINARY_NAME=caffeinate-mcp
VERSION=$(shell git describe --tags --always --dirty)
COMMIT=$(shell git rev-parse --short HEAD)
DATE=$(shell date -u +"%Y-%m-%dT%H:%M:%SZ")
LDFLAGS=-ldflags "-s -w -X main.version=$(VERSION) -X main.commit=$(COMMIT) -X main.date=$(DATE)"

# Default target
help:
	@echo "Available targets:"
	@echo "  build          - Build the binary for current OS/arch"
	@echo "  test           - Run tests"
	@echo "  clean          - Remove built binaries"
	@echo "  install        - Install binary to GOPATH/bin"
	@echo "  release-snapshot - Create a snapshot release with GoReleaser"
	@echo "  lint           - Run golangci-lint"
	@echo "  fmt            - Format code with gofmt"

# Build binary
build:
	go build $(LDFLAGS) -o $(BINARY_NAME) .

# Run tests
test:
	go test -v ./...

# Clean built binaries
clean:
	rm -f $(BINARY_NAME)
	rm -rf dist/

# Install binary
install:
	go install $(LDFLAGS) .

# Create snapshot release
release-snapshot:
	goreleaser release --snapshot --clean

# Run linter
lint:
	@which golangci-lint > /dev/null || (echo "golangci-lint not found. Install from https://golangci-lint.run/usage/install/" && exit 1)
	golangci-lint run

# Format code
fmt:
	go fmt ./...

# Run the server (for development)
run:
	go run $(LDFLAGS) .

# Check version
version:
	@go run $(LDFLAGS) . --version