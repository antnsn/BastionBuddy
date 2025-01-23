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
[![Release](https://github.com/antnsn/BastionBuddy/actions/workflows/release.yml/badge.svg)](https://github.com/antnsn/BastionBuddy/actions/workflows/release.yml)
[![Homebrew](https://img.shields.io/badge/homebrew-available-blue)](https://github.com/antnsn/homebrew-bastionbuddy)

A friendly command-line utility that makes Azure Bastion connections easy and interactive.

## Features

- 🔒 Secure SSH connections to Azure VMs
- 🖥️ Remote Desktop (RDP) support
- 🌐 Port tunneling for remote access
- 💾 Save and reuse connection configurations
- 🔄 Automatic Azure CLI login handling
- ⚡ Smart caching for faster resource listing
- 🎯 Interactive menu navigation with search
- 📝 Command history for quick access

## Prerequisites

Before using BastionBuddy, you need to have the following installed:

1. **Azure CLI** (`az`): Required for Azure authentication and operations
   - Install via [Azure CLI installation guide](https://docs.microsoft.com/en-us/cli/azure/install-azure-cli)
   - Quick install commands:
     ```bash
     # macOS
     brew update && brew install azure-cli

     # Windows
     winget install -e --id Microsoft.AzureCLI

     # Ubuntu/Debian
     curl -sL https://aka.ms/InstallAzureCLIDeb | sudo bash
     ```
   - BastionBuddy will automatically prompt you to log in if needed

2. **An Azure account with Bastion service enabled**
   - Ensure you have appropriate permissions to access Azure resources
   - Bastion service must be deployed in your target VNet

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

### Tips for Port Tunnels

When using port tunnels with SSH, you can avoid host key verification issues by using the following SSH command:
```bash
ssh -o UserKnownHostsFile=/dev/null -o StrictHostKeyChecking=no username@localhost -p <local_port>
```

This approach:
- Bypasses host key verification for localhost connections
- Prevents conflicts when connecting to different hosts through the same local port
- Makes it easier to use tools that rely on SSH connections

Example:
```bash
# Start a tunnel to remote port 22 on local port 50021
bastionbuddy tunnel

# Connect using SSH through the tunnel
ssh -o UserKnownHostsFile=/dev/null -o StrictHostKeyChecking=no sysadmin@localhost -p 50021
```

## Configuration Files

BastionBuddy stores your connection configurations in JSON files in the following locations:

### Unix-like Systems (macOS, Linux)
```
~/.config/bastionbuddy/
├── ssh.json      # Saved SSH connections
├── rdp.json      # Saved RDP connections
├── tunnels.json  # Saved port tunnels
└── active.json   # Currently active tunnels
```

### Windows
```
%USERPROFILE%\.config\bastionbuddy\
├── ssh.json      # Saved SSH connections
├── rdp.json      # Saved RDP connections
├── tunnels.json  # Saved port tunnels
└── active.json   # Currently active tunnels
```
(typically `C:\Users\<username>\.config\bastionbuddy\`)

Each configuration file stores connection details such as resource names, subscription IDs, and connection-specific parameters.

### Configuration Parameters

Common parameters across all configuration types:
- `name`: Unique identifier for the configuration
- `subscription_id`: Azure subscription ID
- `resource_id`: Full Azure resource ID
- `resource_name`: Name of the target resource
- `bastion_name`: Name of the Bastion host
- `bastion_resource_group`: Resource group of the Bastion host
- `bastion_subscription_id`: Subscription ID of the Bastion host
- `username`: Username for the connection
- `last_used`: Timestamp of last connection

Additional parameters for tunnels:
- `local_port`: Local port to forward from
- `remote_port`: Remote port to forward to
- `connection_type`: Type of connection ("tunnel")
- `id`: Unique identifier for active tunnels
- `pid`: Process ID of the running tunnel
- `status`: Current status of the tunnel ("running")
- `start_time`: When the tunnel was started

## Commands

BastionBuddy supports the following commands:

### List Saved Configurations
```bash
bastionbuddy list          # List all saved configurations
bastionbuddy list ssh      # List saved SSH configurations
bastionbuddy list rdp      # List saved RDP configurations
bastionbuddy list tunnels  # List saved tunnel configurations
```

### SSH Connections
```bash
bastionbuddy ssh                    # Interactive SSH connection setup
bastionbuddy ssh <config-name>      # Connect using a saved SSH configuration
```

### RDP Connections
```bash
bastionbuddy rdp                    # Interactive RDP connection setup
bastionbuddy rdp <config-name>      # Connect using a saved RDP configuration
```

### Port Tunnels
```bash
bastionbuddy tunnel                    # Interactive tunnel setup
bastionbuddy tunnel <config-name>      # Start a saved tunnel configuration
bastionbuddy tunnel stop <tunnel-id>   # Stop a running tunnel
```

When using any command without parameters, BastionBuddy will guide you through an interactive menu to select resources and configure the connection. Your configurations are automatically saved for future use.

### Managing Active Tunnels

To view and manage active tunnels:
```bash
bastionbuddy
# Select "Manage active tunnels" from the menu
```

Active tunnels are tracked in `active.json` with the following information:
- Tunnel ID
- Local and remote ports
- Resource details
- Process ID (PID)
- Status (running)

You can stop a tunnel either by:
1. Using the interactive menu and selecting "Terminate"
2. Using the command line: `bastionbuddy tunnel stop <tunnel-id>`

The tunnel process will be properly terminated and removed from the active tunnels list.

### Example Usage

1. Create a new SSH connection:
   ```bash
   bastionbuddy ssh
   # Follow the interactive prompts to select:
   # 1. Azure subscription
   # 2. Resource group
   # 3. Virtual machine
   # 4. Username
   ```

2. Reuse a saved SSH configuration:
   ```bash
   bastionbuddy ssh my-dev-server
   ```

3. Create a port tunnel:
   ```bash
   bastionbuddy tunnel
   # Follow the prompts to set up:
   # 1. Local port
   # 2. Remote port
   # 3. Target VM
   ```

4. Connect via RDP:
   ```bash
   bastionbuddy rdp
   # Follow the prompts to select:
   # 1. Target Windows VM
   # 2. Username
   ```

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
