.PHONY: build clean test lint all release

BINARY_NAME=bastionBuddy
BUILD_DIR=builds
VERSION=$(shell git describe --tags --always --dirty)

# Default build for current platform
build:
	mkdir -p $(BUILD_DIR)
	go build -o $(BUILD_DIR)/$(BINARY_NAME) -ldflags="-X 'main.Version=$(VERSION)'" ./cmd/azbastion

# Build for all supported platforms
all:
	mkdir -p $(BUILD_DIR)
	GOOS=darwin GOARCH=amd64 go build -o $(BUILD_DIR)/$(BINARY_NAME)_darwin_amd64 -ldflags="-X 'main.Version=$(VERSION)'" ./cmd/azbastion
	GOOS=darwin GOARCH=arm64 go build -o $(BUILD_DIR)/$(BINARY_NAME)_darwin_arm64 -ldflags="-X 'main.Version=$(VERSION)'" ./cmd/azbastion
	GOOS=linux GOARCH=amd64 go build -o $(BUILD_DIR)/$(BINARY_NAME)_linux_amd64 -ldflags="-X 'main.Version=$(VERSION)'" ./cmd/azbastion
	GOOS=linux GOARCH=arm64 go build -o $(BUILD_DIR)/$(BINARY_NAME)_linux_arm64 -ldflags="-X 'main.Version=$(VERSION)'" ./cmd/azbastion

# Create release archives
release: all
	cd $(BUILD_DIR) && \
	for file in $(BINARY_NAME)_* ; do \
		zip $${file}.zip $$file ; \
		rm $$file ; \
	done

clean:
	rm -rf $(BUILD_DIR)

test:
	go test -v ./...

lint:
	@if ! command -v golangci-lint &> /dev/null; then \
		echo "golangci-lint not found. Run ./scripts/check.sh to install it."; \
		exit 1; \
	fi
	go mod tidy
	golangci-lint run

install: build
	cp $(BUILD_DIR)/$(BINARY_NAME) $(GOPATH)/bin/

.DEFAULT_GOAL := build
