# BastionBuddy

[![Build Status](https://github.com/antnsn/BastionBuddy/actions/workflows/pr-check.yml/badge.svg)](https://github.com/antnsn/BastionBuddy/actions/workflows/pr-check.yml)
[![Release](https://github.com/antnsn/BastionBuddy/actions/workflows/release.yml/badge.svg)](https://github.com/antnsn/BastionBuddy/releases)
[![codecov](https://codecov.io/gh/antnsn/BastionBuddy/branch/main/graph/badge.svg)](https://codecov.io/gh/antnsn/BastionBuddy)

A friendly command-line utility that makes Azure Bastion connections easy and interactive.

## Features

- 🔒 Secure SSH connections to Azure VMs
- 🌐 Port tunneling for remote access
- ⚡ Smart caching for faster resource listing
- 🎯 Interactive menu navigation with search

## Prerequisites

Before using BastionBuddy, you need to have the following installed:

1. **Azure CLI** (`az`): Required for Azure authentication and operations
   - Installation guide: [Install the Azure CLI](https://docs.microsoft.com/en-us/cli/azure/install-azure-cli)
   - Make sure you're logged in with `az login`

That's it! All other dependencies are bundled with the application.

## Installation

1. Download the latest release from the releases page
2. Make it executable:
   ```bash
   chmod +x ./bastionBuddy
   ```
3. Optionally, move it to your PATH:
   ```bash
   sudo mv bastionBuddy /usr/local/bin/
   ```

## Homebrew   

### Add the tap
```bash
brew tap antnsn/bastionbuddy
```
### Install BastionBuddy
```bash
brew install bastionbuddy
```
## Usage

Simply run:
```bash
bastionBuddy
```

The interactive menu will guide you through:
1. Selecting connection type (SSH or Tunnel)
2. Choosing target subscription
3. Selecting target resource
4. Establishing the connection

## Tips

- Use ↑/↓ arrow keys to navigate menus
- Type to search in lists
- Press Enter to select
- Use Ctrl+C to exit at any time

## Development

### Local Testing

Before pushing changes, you can run all checks locally using:
```bash
./scripts/check.sh
```

This will:
1. Run the linter
2. Execute tests with race detection
3. Generate and display coverage report
4. Build for current platform
5. Test build for all supported platforms

A pre-commit hook is also available that runs these checks automatically before each commit.

### Continuous Integration

The following checks run automatically on pull requests:
- Code linting with golangci-lint
- Unit tests with race detection
- Code coverage reporting
- Multi-platform build verification

### Release Process

1. Update version number where needed
2. Create and push a new tag:
   ```bash
   git tag v1.x.x
   git push origin v1.x.x
   ```
3. GitHub Actions will automatically:
   - Build for all platforms
   - Create GitHub release
   - Upload binaries
   - Generate changelog

## Building from Source

1. Clone the repository
2. Run:
   ```bash
   make
   ```
3. The binary will be available in `./bin/bastionBuddy`
