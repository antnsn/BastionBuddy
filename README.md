# BastionBuddy

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

## Building from Source

1. Clone the repository
2. Run:
   ```bash
   make
   ```
3. The binary will be available in `./bin/bastionBuddy`
