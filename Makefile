.PHONY: build run test clean

# Build the application
build:
	go build -o bin/server cmd/server/main.go

# Run the application
run:
	go run cmd/server/main.go

# Run tests
test:
	go test -v ./...

# Clean build artifacts
clean:
	rm -rf bin/

# Install dependencies
deps:
	go mod tidy
	go mod download

# Run linter
lint:
	golangci-lint run

# Generate mock files (if needed)
mocks:
	mockgen -source=internal/services/analyze_service.go -destination=internal/mocks/analyze_service_mock.go
	mockgen -source=internal/services/generate_service.go -destination=internal/mocks/generate_service_mock.go
	mockgen -source=internal/services/pdf_service.go -destination=internal/mocks/pdf_service_mock.go

# Default target
all: clean build test 