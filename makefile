.PHONY: build clean test lint

BINARY_NAME=bastionBuddy
BUILD_DIR=bin

build:
	go build -o $(BUILD_DIR)/$(BINARY_NAME) ./cmd/azbastion

clean:
	rm -rf $(BUILD_DIR)

test:
	go test -v ./...

lint:
	golangci-lint run

install: build
	cp $(BUILD_DIR)/$(BINARY_NAME) $(GOPATH)/bin/

.DEFAULT_GOAL := build
