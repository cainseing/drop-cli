BINARY_NAME=drop
VERSION=0.5.0-beta
LDFLAGS=-ldflags "-X github.com/cainseing/drop-cli/internal/config.Version=v$(VERSION)"

.PHONY: all build install clean test install-hooks publish

all: build

# Build for the current architecture
build:
	@echo "🛠️  Building $(BINARY_NAME)..."
	go build -o $(BINARY_NAME) ./cmd/drop

# Install to your system path (Standard for Linux/macOS)
install: build
	@echo "🚀 Installing to /usr/local/bin..."
	@sudo mv $(BINARY_NAME) /usr/local/bin/$(BINARY_NAME)
	@echo "✅ Done! Try running: drop --help"

# Build for multiple platforms (Cross-Compilation)
release: test
	@echo "🌎 Building for multiple platforms..."
	GOOS=linux GOARCH=amd64 go build $(LDFLAGS) -o bin/$(BINARY_NAME)-linux-amd64 ./cmd/drop
	GOOS=linux GOARCH=arm64 go build $(LDFLAGS) -o bin/$(BINARY_NAME)-linux-arm64 ./cmd/drop
	GOOS=darwin GOARCH=arm64 go build $(LDFLAGS) -o bin/$(BINARY_NAME)-darwin-arm64 ./cmd/drop
	GOOS=darwin GOARCH=amd64 go build $(LDFLAGS) -o bin/$(BINARY_NAME)-darwin-amd64 ./cmd/drop
# 	GOOS=windows GOARCH=amd64 go build -o bin/$(BINARY_NAME)-windows-amd64.exe .
# 	GOOS=windows GOARCH=arm64 go build -o bin/$(BINARY_NAME)-windows-arm64.exe .

publish:
	@echo "Tagging and releasing v$(VERSION)..."
	git tag v$(VERSION)
	git push origin v$(VERSION)

git-hooks:
	@git config core.hooksPath .githooks
	@echo "Git hooks configured. Tag a release with: git tag v<version> && git push origin v<version>"

clean:
	@echo "Cleaning up..."
	@rm -f $(BINARY_NAME)
	@rm -rf bin/

test:
	@echo "Running tests..."
	go test ./...

test-coverage:
	@echo "Running tests with coverage..."
	go test -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out -o coverage.html
	@echo "📊 Coverage report generated: coverage.html"

lint:
	@echo "Running linter..."
	golangci-lint run

sec:
	@echo "Running gosec..."
	gosec ./...