.PHONY: all build test coverage clean

# Default target
all: test build

# Build the application
build:
	go build -o review-extractor ./cmd

# Run all tests
test:
	go test -v ./...

# Run tests with coverage
coverage:
	go test -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out

# Clean build artifacts
clean:
	rm -f review-extractor
	rm -f coverage.out

# Update dependencies
deps:
	go mod tidy
	go mod verify

# Run linter
lint:
	golangci-lint run

# Run all checks (tests, lint, coverage)
check: test lint coverage 