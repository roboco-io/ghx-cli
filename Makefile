# ghp-cli Makefile

# Go parameters
GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get
GOMOD=$(GOCMD) mod
GOFMT=gofmt
GOLINT=golangci-lint
BINARY_NAME=ghp
BINARY_PATH=bin/$(BINARY_NAME)

# Version info
VERSION ?= $(shell git describe --tags --always --dirty)
COMMIT ?= $(shell git rev-parse --short HEAD)
BUILD_TIME ?= $(shell date -u '+%Y-%m-%d_%H:%M:%S')
LDFLAGS := -ldflags "-X main.Version=$(VERSION) -X main.Commit=$(COMMIT) -X main.BuildTime=$(BUILD_TIME)"

# Test parameters
TEST_FLAGS ?= -v -race -coverprofile=coverage.out
TEST_TIMEOUT ?= 10m

.PHONY: all build clean test test-unit test-integration test-e2e test-e2e-write coverage fmt lint install deps help

## help: Show this help message
help:
	@echo 'Usage:'
	@sed -n 's/^##//p' ${MAKEFILE_LIST} | column -t -s ':' | sed -e 's/^/ /'

## all: Run tests and build
all: test build

## build: Build the binary
build:
	@echo "Building $(BINARY_NAME)..."
	@mkdir -p bin
	$(GOBUILD) $(LDFLAGS) -o $(BINARY_PATH) ./cmd/ghp

## clean: Remove build artifacts
clean:
	@echo "Cleaning..."
	$(GOCLEAN)
	@rm -rf bin/ coverage.out coverage.html dist/

## test: Run tests
test:
	@echo "Running tests..."
	$(GOTEST) $(TEST_FLAGS) -timeout $(TEST_TIMEOUT) ./...

## test-unit: Run unit tests only
test-unit:
	@echo "Running unit tests..."
	$(GOTEST) $(TEST_FLAGS) -short -timeout $(TEST_TIMEOUT) ./...

## test-integration: Run integration tests
test-integration:
	@echo "Running integration tests..."
	$(GOTEST) $(TEST_FLAGS) -run Integration -timeout $(TEST_TIMEOUT) ./test/...

## test-e2e: Run E2E tests (requires GITHUB_TOKEN or gh auth)
test-e2e: build
	@echo "Running E2E tests..."
	@echo "Note: Set GHP_TEST_REPO=owner/repo to use a different test repository"
	$(GOTEST) -v -timeout $(TEST_TIMEOUT) ./test/e2e/...

## test-e2e-write: Run E2E tests including write operations (creates/deletes discussions)
test-e2e-write: build
	@echo "Running E2E tests with write operations..."
	@echo "WARNING: This will create and delete discussions in the test repository"
	GHP_E2E_WRITE_TESTS=1 $(GOTEST) -v -timeout $(TEST_TIMEOUT) ./test/e2e/...

## test-all: Run all tests (unit, integration, e2e)
test-all: test test-e2e
	@echo "All tests completed"

## coverage: Generate test coverage report
coverage: test
	@echo "Generating coverage report..."
	@go tool cover -html=coverage.out -o coverage.html
	@echo "Coverage report generated: coverage.html"

## coverage-view: View coverage in browser
coverage-view: coverage
	@open coverage.html

## fmt: Format code
fmt:
	@echo "Formatting code..."
	$(GOFMT) -s -w .
	@go mod tidy

## lint: Run linter
lint:
	@echo "Running linter..."
	@if ! which $(GOLINT) > /dev/null; then \
		echo "Installing golangci-lint..."; \
		curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(shell go env GOPATH)/bin v1.62.2; \
	fi
	$(GOLINT) run ./...

## install: Install the binary
install: build
	@echo "Installing $(BINARY_NAME)..."
	@cp $(BINARY_PATH) $(GOPATH)/bin/$(BINARY_NAME)

## deps: Download dependencies
deps:
	@echo "Downloading dependencies..."
	$(GOMOD) download
	$(GOMOD) tidy

## deps-update: Update dependencies
deps-update:
	@echo "Updating dependencies..."
	$(GOGET) -u ./...
	$(GOMOD) tidy

## setup: Setup development environment
setup: deps install-tools
	@echo "Setting up pre-commit hooks..."
	@npm install --save-dev husky
	@npx husky init
	@echo "Development environment ready!"

## install-tools: Install development tools
install-tools:
	@echo "Installing development tools..."
	@go install github.com/golangci/golangci-lint/cmd/golangci-lint@v1.62.2
	@go install github.com/goreleaser/goreleaser/v2@v2.5.1
	@go install golang.org/x/tools/cmd/goimports@latest
	@go install github.com/securego/gosec/v2/cmd/gosec@latest
	@go install github.com/mgechev/revive@latest

## run: Run the application
run:
	@echo "Running $(BINARY_NAME)..."
	$(GOCMD) run ./cmd/ghp

## docker-build: Build Docker image
docker-build:
	@echo "Building Docker image..."
	@docker build -t ghp-cli:$(VERSION) .

## release: Create a new release
release:
	@echo "Creating release..."
	@goreleaser release --clean

## release-snapshot: Create a snapshot release
release-snapshot:
	@echo "Creating snapshot release..."
	@goreleaser release --snapshot --clean