.PHONY: build run clean install test

# Build the binary
build:
	go build -o bin/cp-discovery ./cmd/cp-discovery

# Build with all optimizations
build-release:
	CGO_ENABLED=1 go build -ldflags="-s -w" -o bin/cp-discovery ./cmd/cp-discovery

# Run the tool
run: build
	./bin/cp-discovery

# Run with custom config
run-config: build
	./bin/cp-discovery -config $(CONFIG)

# Install dependencies
install:
	go mod download
	go mod tidy

# Clean build artifacts
clean:
	rm -rf bin/
	rm -f discovery-report.*

# Run tests (when added)
test:
	go test -v ./...

# Format code
fmt:
	go fmt ./...

# Vet code
vet:
	go vet ./...

# Run all checks
check: fmt vet

# Install to system (optional)
install-bin: build-release
	sudo cp bin/cp-discovery /usr/local/bin/

# Help
help:
	@echo "Available targets:"
	@echo "  build         - Build the binary"
	@echo "  build-release - Build optimized binary"
	@echo "  run           - Build and run with default config"
	@echo "  run-config    - Build and run with custom config (make run-config CONFIG=path/to/config.yaml)"
	@echo "  install       - Install Go dependencies"
	@echo "  clean         - Remove build artifacts"
	@echo "  test          - Run tests"
	@echo "  fmt           - Format code"
	@echo "  vet           - Run go vet"
	@echo "  check         - Run fmt and vet"
	@echo "  install-bin   - Install binary to /usr/local/bin"
	@echo "  help          - Show this help message"
