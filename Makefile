# Makefile for quanto-magis project

# Go parameters
GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get
GOMOD=$(GOCMD) mod
GOFMT=$(GOCMD) fmt
GOVET=$(GOCMD) vet

# Binary name and paths
BINARY_NAME=quanto
BINARY_UNIX=$(BINARY_NAME)_unix
BINARY_DARWIN=$(BINARY_NAME)_darwin
BINARY_WINDOWS=$(BINARY_NAME).exe

# Build directory
BUILD_DIR=build
CMD_DIR=./cmd/quanto

# Version information
VERSION?=$(shell git describe --tags --always --dirty 2>/dev/null || echo "dev")
COMMIT?=$(shell git rev-parse --short HEAD 2>/dev/null || echo "unknown")
BUILD_TIME?=$(shell date -u '+%Y-%m-%d_%H:%M:%S')

# Linker flags
LDFLAGS=-ldflags "-X main.Version=$(VERSION) -X main.Commit=$(COMMIT) -X main.BuildTime=$(BUILD_TIME)"

# Test parameters
TEST_FLAGS=-v -race -coverprofile=coverage.out -covermode=atomic
TEST_TIMEOUT=10m

# Colors for output
CYAN=\033[0;36m
GREEN=\033[0;32m
RED=\033[0;31m
YELLOW=\033[0;33m
NC=\033[0m # No Color

.PHONY: all build clean test coverage run install uninstall fmt vet lint deps help \
        build-linux build-darwin build-windows build-all docker-build docker-run \
        test-unit test-integration test-all bench

.DEFAULT_GOAL := help

## help: Display this help message
help:
	@echo "$(CYAN)Available targets:$(NC)"
	@grep -E '^##[[:space:]]' $(MAKEFILE_LIST) | sed 's/^##[[:space:]]//' | awk 'BEGIN {FS = ":"}; {printf "  $(GREEN)%-20s$(NC) %s\n", $$1, $$2}'

##@ Build

## build: Build the binary for current platform
build:
	@echo "$(CYAN)Building $(BINARY_NAME)...$(NC)"
	@mkdir -p $(BUILD_DIR)
	$(GOBUILD) $(LDFLAGS) -o $(BUILD_DIR)/$(BINARY_NAME) $(CMD_DIR)
	@echo "$(GREEN)Build complete: $(BUILD_DIR)/$(BINARY_NAME)$(NC)"

## build-linux: Build for Linux
build-linux:
	@echo "$(CYAN)Building for Linux...$(NC)"
	@mkdir -p $(BUILD_DIR)
	GOOS=linux GOARCH=amd64 $(GOBUILD) $(LDFLAGS) -o $(BUILD_DIR)/$(BINARY_UNIX) $(CMD_DIR)
	@echo "$(GREEN)Build complete: $(BUILD_DIR)/$(BINARY_UNIX)$(NC)"

## build-darwin: Build for macOS
build-darwin:
	@echo "$(CYAN)Building for macOS...$(NC)"
	@mkdir -p $(BUILD_DIR)
	GOOS=darwin GOARCH=amd64 $(GOBUILD) $(LDFLAGS) -o $(BUILD_DIR)/$(BINARY_DARWIN) $(CMD_DIR)
	@echo "$(GREEN)Build complete: $(BUILD_DIR)/$(BINARY_DARWIN)$(NC)"

## build-windows: Build for Windows
build-windows:
	@echo "$(CYAN)Building for Windows...$(NC)"
	@mkdir -p $(BUILD_DIR)
	GOOS=windows GOARCH=amd64 $(GOBUILD) $(LDFLAGS) -o $(BUILD_DIR)/$(BINARY_WINDOWS) $(CMD_DIR)
	@echo "$(GREEN)Build complete: $(BUILD_DIR)/$(BINARY_WINDOWS)$(NC)"

## build-all: Build for all platforms
build-all: build-linux build-darwin build-windows
	@echo "$(GREEN)All platform builds complete$(NC)"

##@ Run

## run: Run the application
run: build
	@echo "$(CYAN)Running $(BINARY_NAME)...$(NC)"
	@$(BUILD_DIR)/$(BINARY_NAME)

## run-dev: Run the application without building (go run)
run-dev:
	@echo "$(CYAN)Running in development mode...$(NC)"
	$(GOCMD) run $(CMD_DIR)/main.go

##@ Testing

## test: Run all tests
test:
	@echo "$(CYAN)Running tests...$(NC)"
	$(GOTEST) $(TEST_FLAGS) -timeout $(TEST_TIMEOUT) ./...
	@echo "$(GREEN)Tests complete$(NC)"

## test-unit: Run unit tests only
test-unit:
	@echo "$(CYAN)Running unit tests...$(NC)"
	$(GOTEST) -v -race -short -timeout $(TEST_TIMEOUT) ./...
	@echo "$(GREEN)Unit tests complete$(NC)"

## test-coverage: Run tests with coverage report
test-coverage:
	@echo "$(CYAN)Running tests with coverage...$(NC)"
	$(GOTEST) $(TEST_FLAGS) -timeout $(TEST_TIMEOUT) ./...
	$(GOCMD) tool cover -html=coverage.out -o coverage.html
	@echo "$(GREEN)Coverage report generated: coverage.html$(NC)"

## test-integration: Run integration tests
test-integration:
	@echo "$(CYAN)Running integration tests...$(NC)"
	$(GOTEST) -v -race -run Integration -timeout $(TEST_TIMEOUT) ./...
	@echo "$(GREEN)Integration tests complete$(NC)"

## bench: Run benchmarks
bench:
	@echo "$(CYAN)Running benchmarks...$(NC)"
	$(GOTEST) -bench=. -benchmem ./...
	@echo "$(GREEN)Benchmarks complete$(NC)"

## coverage: Show test coverage
coverage:
	@echo "$(CYAN)Calculating coverage...$(NC)"
	$(GOTEST) -coverprofile=coverage.out ./...
	$(GOCMD) tool cover -func=coverage.out
	@echo "$(GREEN)Total coverage:$(NC)"
	@$(GOCMD) tool cover -func=coverage.out | grep total | awk '{print $$3}'

##@ Code Quality

## fmt: Format code
fmt:
	@echo "$(CYAN)Formatting code...$(NC)"
	$(GOFMT) ./...
	@echo "$(GREEN)Formatting complete$(NC)"

## vet: Run go vet
vet:
	@echo "$(CYAN)Running go vet...$(NC)"
	$(GOVET) ./...
	@echo "$(GREEN)Vet complete$(NC)"

## lint: Run golangci-lint
lint:
	@echo "$(CYAN)Running linter...$(NC)"
	@which golangci-lint > /dev/null 2>&1 || (echo "$(RED)golangci-lint not installed. Run 'make install-tools'$(NC)" && exit 1)
	golangci-lint run --config .golangci.yml
	@echo "$(GREEN)Linting complete$(NC)"

## check: Run fmt, vet, and lint
check: fmt vet lint
	@echo "$(GREEN)All checks passed$(NC)"

##@ Dependencies

## deps: Download dependencies
deps:
	@echo "$(CYAN)Downloading dependencies...$(NC)"
	$(GOMOD) download
	@echo "$(GREEN)Dependencies downloaded$(NC)"

## deps-tidy: Tidy dependencies
deps-tidy:
	@echo "$(CYAN)Tidying dependencies...$(NC)"
	$(GOMOD) tidy
	@echo "$(GREEN)Dependencies tidied$(NC)"

## deps-verify: Verify dependencies
deps-verify:
	@echo "$(CYAN)Verifying dependencies...$(NC)"
	$(GOMOD) verify
	@echo "$(GREEN)Dependencies verified$(NC)"

## deps-update: Update all dependencies
deps-update:
	@echo "$(CYAN)Updating dependencies...$(NC)"
	$(GOGET) -u ./...
	$(GOMOD) tidy
	@echo "$(GREEN)Dependencies updated$(NC)"

##@ Installation

## install: Install the binary to GOPATH/bin
install:
	@echo "$(CYAN)Installing $(BINARY_NAME)...$(NC)"
	$(GOCMD) install $(LDFLAGS) $(CMD_DIR)
	@echo "$(GREEN)Installation complete$(NC)"

## uninstall: Remove the binary from GOPATH/bin
uninstall:
	@echo "$(CYAN)Uninstalling $(BINARY_NAME)...$(NC)"
	@rm -f $(GOPATH)/bin/$(BINARY_NAME)
	@echo "$(GREEN)Uninstallation complete$(NC)"

## install-tools: Install development tools
install-tools:
	@echo "$(CYAN)Installing development tools...$(NC)"
	@which golangci-lint > /dev/null 2>&1 || \
		(echo "Installing golangci-lint..." && \
		curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $$(go env GOPATH)/bin)
	@echo "$(GREEN)Development tools installed$(NC)"

##@ Cleanup

## clean: Remove build artifacts and cache
clean:
	@echo "$(CYAN)Cleaning...$(NC)"
	$(GOCLEAN)
	@rm -rf $(BUILD_DIR)
	@rm -f coverage.out coverage.html
	@echo "$(GREEN)Clean complete$(NC)"

## clean-all: Remove all generated files including dependencies
clean-all: clean
	@echo "$(CYAN)Removing dependencies cache...$(NC)"
	$(GOCLEAN) -cache -testcache -modcache
	@echo "$(GREEN)All clean$(NC)"

##@ Docker (if needed)

## docker-build: Build Docker image
docker-build:
	@echo "$(CYAN)Building Docker image...$(NC)"
	docker build -t $(BINARY_NAME):$(VERSION) .
	@echo "$(GREEN)Docker image built: $(BINARY_NAME):$(VERSION)$(NC)"

## docker-run: Run Docker container
docker-run:
	@echo "$(CYAN)Running Docker container...$(NC)"
	docker run --rm -it $(BINARY_NAME):$(VERSION)

##@ Info

## version: Display version information
version:
	@echo "Version: $(VERSION)"
	@echo "Commit: $(COMMIT)"
	@echo "Build Time: $(BUILD_TIME)"

## info: Display project information
info:
	@echo "$(CYAN)Project Information:$(NC)"
	@echo "  Binary Name: $(BINARY_NAME)"
	@echo "  Build Dir:   $(BUILD_DIR)"
	@echo "  Version:     $(VERSION)"
	@echo "  Commit:      $(COMMIT)"
	@echo "  Go Version:  $$($(GOCMD) version)"

##@ All-in-one

## all: Clean, format, lint, test, and build
all: clean fmt vet lint test build
	@echo "$(GREEN)Build pipeline complete!$(NC)"

## ci: Run CI pipeline (deps, check, test)
ci: deps check test
	@echo "$(GREEN)CI pipeline complete!$(NC)"
