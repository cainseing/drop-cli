BINARY_NAME=drop
VERSION=0.1.0-beta

.PHONY: all build install clean test

all: build

# Build for the current architecture
build:
	@echo "🛠️  Building $(BINARY_NAME)..."
	go build -o $(BINARY_NAME) .

# Install to your system path (Standard for Linux/macOS)
install: build
	@echo "🚀 Installing to /usr/local/bin..."
	@sudo mv $(BINARY_NAME) /usr/local/bin/$(BINARY_NAME)
	@echo "✅ Done! Try running: drop --help"

# Build for multiple platforms (Cross-Compilation)
release:
	@echo "🌎 Building for multiple platforms..."
	GOOS=linux GOARCH=amd64 go build -o bin/$(BINARY_NAME)-linux-amd64 .
	GOOS=linux GOARCH=arm64 go build -o bin/$(BINARY_NAME)-linux-arm64 .
	GOOS=darwin GOARCH=arm64 go build -o bin/$(BINARY_NAME)-darwin-arm64 .
	GOOS=darwin GOARCH=amd64 go build -o bin/$(BINARY_NAME)-darwin-amd64 .
# 	GOOS=windows GOARCH=amd64 go build -o bin/$(BINARY_NAME)-windows-amd64.exe .
# 	GOOS=windows GOARCH=arm64 go build -o bin/$(BINARY_NAME)-windows-arm64.exe .

clean:
	@echo "🧹 Cleaning up..."
	@rm -f $(BINARY_NAME)
	@rm -rf bin/

test:
	@echo "🧪 Running tests..."
	go test ./...