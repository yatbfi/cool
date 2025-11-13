#!make
BINARY := cool
GOARCH := amd64
GOOS   ?= $(shell go env GOOS)
VERSION ?= development

COMMIT := $(shell git rev-parse --short HEAD)
BRANCH := $(shell git rev-parse --abbrev-ref HEAD)
DATE   := $(shell date -u +"%Y-%m-%dT%H:%M:%SZ")

LDFLAGS = -ldflags "-X 'main.Version=$(VERSION)' -X 'main.Commit=$(COMMIT)' -X 'main.Branch=$(BRANCH)' -X 'main.BuildDate=$(DATE)'"

# Default goal
.DEFAULT_GOAL := help

# ------------------------------------------------------------------------------

.PHONY: lint
# Run static analysis with golangci-lint
lint:
	@golangci-lint run

.PHONY: build
# Build current OS binary
build:
	@echo "🚀 Building for $(GOOS)/$(GOARCH)..."
	@go build $(LDFLAGS) -o bin/$(BINARY) .
	@echo "✅ Built: bin/$(BINARY)"

.PHONY: install
# Install binary to GOBIN
install:
	@echo "📦 Installing $(BINARY) to $$GOBIN..."
	@go install $(LDFLAGS) .

.PHONY: linux
# Cross-compile for Linux
linux:
	@echo "🐧 Building for linux/$(GOARCH)..."
	GOOS=linux GOARCH=$(GOARCH) go build $(LDFLAGS) -o bin/$(BINARY)-linux-$(GOARCH) .
	@echo "✅ Built: bin/$(BINARY)-linux-$(GOARCH)"

.PHONY: darwin
# Cross-compile for macOS
darwin:
	@echo "🍎 Building for darwin/$(GOARCH)..."
	GOOS=darwin GOARCH=$(GOARCH) go build $(LDFLAGS) -o bin/$(BINARY)-darwin-$(GOARCH) .
	@echo "✅ Built: bin/$(BINARY)-darwin-$(GOARCH)"

.PHONY: windows
# Cross-compile for Windows
windows:
	@echo "🪟 Building for windows/$(GOARCH)..."
	GOOS=windows GOARCH=$(GOARCH) go build $(LDFLAGS) -o bin/$(BINARY)-windows-$(GOARCH).exe .
	@echo "✅ Built: bin/$(BINARY)-windows-$(GOARCH).exe"

.PHONY: all
# Build binaries for all platforms
all: clean linux darwin windows
	@echo "🎉 All platform builds completed!"

.PHONY: clean
# Remove build artifacts
clean:
	@echo "🧹 Cleaning build artifacts..."
	@rm -rf bin
	@echo "✅ Clean done"

# ------------------------------------------------------------------------------

help:
	@echo ''
	@echo 'Usage:'
	@echo '  make [target]'
	@echo ''
	@echo 'Targets:'
	@awk '/^[a-zA-Z\-\_0-9]+:/ { \
	helpMessage = match(lastLine, /^# (.*)/); \
		if (helpMessage) { \
			helpCommand = substr($$1, 0, index($$1, ":")-1); \
			helpMessage = substr(lastLine, RSTART + 2, RLENGTH); \
			printf "\033[36m%-20s\033[0m %s\n", helpCommand,helpMessage; \
		} \
	} \
	{ lastLine = $$0 }' $(MAKEFILE_LIST)