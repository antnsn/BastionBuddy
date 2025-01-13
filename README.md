  ```
  ██████╗  █████╗ ███████╗████████╗██╗ ██████╗ ███╗   ██╗
  ██╔══██╗██╔══██╗██╔════╝╚══██╔══╝██║██╔═══██╗████╗  ██║
  ██████╔╝███████║███████╗   ██║   ██║██║   ██║██╔██╗ ██║
  ██╔══██╗██╔══██║╚════██║   ██║   ██║██║   ██║██║╚██╗██║
  ██████╔╝██║  ██║███████║   ██║   ██║╚██████╔╝██║ ╚████║
  ╚═════╝ ╚═╝  ╚═╝╚══════╝   ╚═╝   ╚═╝ ╚═════╝ ╚═╝  ╚═══╝

         ██████╗ ██╗   ██╗██████╗ ██████╗ ██╗   ██╗
         ██╔══██╗██║   ██║██╔══██╗██╔══██╗╚██╗ ██╔╝
         ██████╔╝██║   ██║██║  ██║██║  ██║ ╚████╔╝
         ██╔══██╗██║   ██║██║  ██║██║  ██║  ╚██╔╝
         ██████╔╝╚██████╔╝██████╔╝██████╔╝   ██║
         ╚═════╝  ╚═════╝ ╚═════╝ ╚═════╝    ╚═╝
```

[![Build Status](https://github.com/antnsn/BastionBuddy/actions/workflows/pr-check.yml/badge.svg)](https://github.com/antnsn/BastionBuddy/actions/workflows/pr-check.yml)
[![Release](https://github.com/antnsn/BastionBuddy/actions/workflows/release.yml/badge.svg)](https://github.com/antnsn/BastionBuddy/releases)
[![Homebrew](https://img.shields.io/badge/homebrew-available-blue)](https://github.com/antnsn/homebrew-bastionbuddy)

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

### Homebrew (Recommended)
```bash
# Add the tap
brew tap antnsn/bastionbuddy

# Install BastionBuddy
brew install bastionbuddy
```

### Manual Installation

1. Download the appropriate version for your platform from the [releases page](https://github.com/antnsn/BastionBuddy/releases)
   - macOS ARM64 (M1/M2): `bastionbuddy_darwin_arm64.tar.gz`
   - macOS Intel: `bastionbuddy_darwin_amd64.tar.gz`
   - Linux ARM64: `bastionbuddy_linux_arm64.tar.gz`
   - Linux x86_64: `bastionbuddy_linux_amd64.tar.gz`

2. Extract the archive:
   ```bash
   tar xzf bastionbuddy_*_*.tar.gz
   ```

3. Make it executable and move to your PATH:
   ```bash
   chmod +x bastionbuddy
   sudo mv bastionbuddy /usr/local/bin/
   ```

## Usage

Simply run:
```bash
bastionbuddy
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

1. Create and push a new tag:
   ```bash
   git tag v1.x.x
   git push origin v1.x.x
   ```

2. GitHub Actions will automatically:
   - Build binaries for all platforms
   - Create a GitHub release
   - Upload the binaries
   - Generate a changelog
   - Update the Homebrew formula

## Building from Source

1. Clone the repository:
   ```bash
   git clone https://github.com/antnsn/BastionBuddy.git
   cd BastionBuddy
   ```

2. Build for your platform:
   ```bash
   make build
   ```

   Or build for all platforms:
   ```bash
   make all
   ```

The binaries will be in the `builds` directory.
