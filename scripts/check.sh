#!/bin/bash

# ANSI color codes
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Function to print section headers
print_header() {
    echo -e "\n${BLUE}=== $1 ===${NC}\n"
}

# Function to check command result
check_result() {
    if [ $? -eq 0 ]; then
        echo -e "${GREEN}✓ $1 passed${NC}"
        return 0
    else
        echo -e "${RED}✗ $1 failed${NC}"
        return 1
    fi
}

# Function to install golangci-lint
install_golangci_lint() {
    print_header "Installing golangci-lint"
    
    # Check the operating system
    if [[ "$OSTYPE" == "darwin"* ]]; then
        # macOS
        if command -v brew &> /dev/null; then
            echo "Installing via Homebrew..."
            brew install golangci-lint
        else
            echo "Installing via Go..."
            go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
        fi
    elif [[ "$OSTYPE" == "linux-gnu"* ]]; then
        # Linux
        curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(go env GOPATH)/bin v1.55.2
    else
        # Other OS - try Go installation
        go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
    fi

    # Verify installation
    if ! command -v golangci-lint &> /dev/null; then
        echo -e "${RED}Failed to install golangci-lint${NC}"
        return 1
    fi
}

# Initialize error counter
ERRORS=0

# Start checks
print_header "Starting local checks"

# Check if golangci-lint is installed
if ! command -v golangci-lint &> /dev/null; then
    echo -e "${YELLOW}Warning: golangci-lint is not installed.${NC}"
    install_golangci_lint || exit 1
fi

# Clean any previous builds
print_header "Cleaning previous builds"
make clean
check_result "Clean" || ((ERRORS++))

# Ensure dependencies are tidy
print_header "Checking dependencies"
go mod tidy
check_result "Go mod tidy" || ((ERRORS++))

# Format Go files
print_header "Formatting Go files"
find . -name "*.go" -exec gofmt -w {} +
check_result "Go formatting" || ((ERRORS++))

# Run linting
print_header "Running linter"
golangci-lint run
check_result "Linting" || ((ERRORS++))

# Run tests with race detection and coverage
print_header "Running tests with race detection and coverage"
go test -race -coverprofile=coverage.txt -covermode=atomic ./...
check_result "Tests" || ((ERRORS++))

# Show coverage report
if [ -f coverage.txt ]; then
    print_header "Coverage Report"
    go tool cover -func=coverage.txt
    rm coverage.txt
fi

# Build for current platform
print_header "Building for current platform"
make build
check_result "Build" || ((ERRORS++))

# Test build for all platforms
print_header "Testing build for all platforms"
make all
check_result "Multi-platform build" || ((ERRORS++))

# Final summary
print_header "Summary"
if [ $ERRORS -eq 0 ]; then
    echo -e "${GREEN}All checks passed successfully!${NC}"
    echo -e "\nYou can now commit and push your changes."
else
    echo -e "${RED}${ERRORS} check(s) failed!${NC}"
    echo -e "\nPlease fix the issues before committing."
    exit 1
fi
