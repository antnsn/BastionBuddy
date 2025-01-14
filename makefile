.PHONY: build clean test lint all release cleanup

BINARY_NAME=bastionbuddy
BUILD_DIR=builds
VERSION=$(shell git describe --tags --always --dirty)
PLATFORMS=darwin/amd64 darwin/arm64 linux/amd64 linux/arm64 windows/amd64 windows/arm64

# Default build for current platform
build:
	mkdir -p $(BUILD_DIR)
	go build -o $(BUILD_DIR)/$(BINARY_NAME) -ldflags="-X 'main.Version=$(VERSION)' -X 'github.com/antnsn/BastionBuddy/internal/welcome.Version=$(VERSION)'" ./cmd/azbastion

# Build for all supported platforms
all:
	@for platform in $(PLATFORMS); do \
		GOOS=$${platform%/*} ; \
		GOARCH=$${platform#*/} ; \
		OUTPUT_DIR=$(BUILD_DIR)/$${GOOS}_$${GOARCH} ; \
		echo "Building for $${GOOS}/$${GOARCH}..." ; \
		mkdir -p $${OUTPUT_DIR} ; \
		if [ "$${GOOS}" = "windows" ]; then \
			GOOS=$${GOOS} GOARCH=$${GOARCH} go build -o $${OUTPUT_DIR}/$(BINARY_NAME).exe -ldflags="-X 'main.Version=$(VERSION)' -X 'github.com/antnsn/BastionBuddy/internal/welcome.Version=$(VERSION)'" ./cmd/azbastion ; \
		else \
			GOOS=$${GOOS} GOARCH=$${GOARCH} go build -o $${OUTPUT_DIR}/$(BINARY_NAME) -ldflags="-X 'main.Version=$(VERSION)' -X 'github.com/antnsn/BastionBuddy/internal/welcome.Version=$(VERSION)'" ./cmd/azbastion ; \
		fi \
	done

# Create release archives
release: all
	@cd $(BUILD_DIR) && \
	for platform in $(PLATFORMS); do \
		GOOS=$${platform%/*} ; \
		GOARCH=$${platform#*/} ; \
		echo "Creating archive for $${GOOS}/$${GOARCH}..." ; \
		if [ "$${GOOS}" = "windows" ]; then \
			zip -j $(BINARY_NAME)_$${GOOS}_$${GOARCH}.zip $${GOOS}_$${GOARCH}/$(BINARY_NAME).exe ; \
		else \
			tar czf $(BINARY_NAME)_$${GOOS}_$${GOARCH}.tar.gz -C $${GOOS}_$${GOARCH} $(BINARY_NAME) ; \
		fi \
	done

# Clean up build directories
cleanup:
	@cd $(BUILD_DIR) && rm -rf darwin_* linux_* windows_*

clean:
	rm -rf $(BUILD_DIR)

test:
	go test -v ./...

lint:
	@if [ -f "/home/runner/golangci-lint-1.63.4-linux-amd64/golangci-lint" ]; then \
		/home/runner/golangci-lint-1.63.4-linux-amd64/golangci-lint run; \
	elif command -v golangci-lint &> /dev/null; then \
		golangci-lint run; \
	else \
		echo "golangci-lint not found. Run ./scripts/check.sh to install it."; \
		exit 1; \
	fi
	go mod tidy

install: build
	cp $(BUILD_DIR)/$(BINARY_NAME) $(GOPATH)/bin/

.DEFAULT_GOAL := build
