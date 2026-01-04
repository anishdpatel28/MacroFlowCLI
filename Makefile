.PHONY: build install clean test fmt lint

# Build the binary
build:
	go build -o bin/macro main.go

# Build for all platforms
build-all:
	GOOS=darwin GOARCH=amd64 go build -o bin/macro-darwin-amd64 main.go
	GOOS=darwin GOARCH=arm64 go build -o bin/macro-darwin-arm64 main.go
	GOOS=linux GOARCH=amd64 go build -o bin/macro-linux-amd64 main.go
	GOOS=windows GOARCH=amd64 go build -o bin/macro-windows-amd64.exe main.go

# Install locally
install: build
	sudo cp bin/macro /usr/local/bin/macro
	@echo "✓ Installed macro to /usr/local/bin/macro"

# Clean build artifacts
clean:
	rm -rf bin/
	go clean

# Run tests
test:
	go test -v ./...

# Format code
fmt:
	go fmt ./...

# Run linter
lint:
	golangci-lint run || go vet ./...

# Get dependencies
deps:
	go mod download
	go mod tidy

# Development build and install
dev: build
	sudo cp bin/macro /usr/local/bin/macro
	@echo "✓ Development build installed"

# Run
run:
	go run main.go
