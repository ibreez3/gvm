.PHONY: build clean test release

BINARY_NAME=gvm
BUILD_DIR=bin

build:
	@echo "Building $(BINARY_NAME)..."
	@go build -o $(BINARY_NAME) ./cmd/gvm

clean:
	@echo "Cleaning..."
	@rm -f $(BINARY_NAME)
	@rm -rf dist/

test:
	@echo "Running tests..."
	@go test ./...

release:
	@./scripts/release.sh
