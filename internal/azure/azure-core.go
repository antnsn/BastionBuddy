// Package azure provides functionality for interacting with Azure resources
// through the Azure SDK, specifically for Bastion connections.
package azure

import (
	"context"
	"fmt"
	"os"
	"strings"

	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/subscription/armsubscription"
	"github.com/antnsn/BastionBuddy/internal/utils"
)

func init() {
	ctx = context.Background()
}

// CheckDependencies verifies that all required Azure dependencies are available.
func CheckDependencies() error {
	var err error
	cred, err = azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		return fmt.Errorf("failed to create credential: %v", err)
	}
	return nil
}

// Cleanup performs any necessary cleanup operations.
func Cleanup() error {
	return nil
}

// ensureAuthenticated ensures the user is logged into Azure.
func ensureAuthenticated() error {
	if cred == nil {
		var err error
		cred, err = azidentity.NewDefaultAzureCredential(nil)
		if err != nil {
			return fmt.Errorf("authentication failed: %v", err)
		}
	}
	return nil
}

// selectSubscription prompts the user to select an Azure subscription.
func selectSubscription(prompt string) error {
	fmt.Printf("Fetching subscriptions for %s...\n", prompt)

	// Create subscription client
	client, err := armsubscription.NewSubscriptionsClient(cred, nil)
	if err != nil {
		return fmt.Errorf("failed to create subscription client: %v", err)
	}

	// List all subscriptions
	pager := client.NewListPager(nil)
	var items []string

	for pager.More() {
		page, err := pager.NextPage(ctx)
		if err != nil {
			return fmt.Errorf("failed to list subscriptions: %v", err)
		}
		for _, sub := range page.Value {
			if sub.DisplayName == nil || sub.SubscriptionID == nil {
				continue
			}
			items = append(items, fmt.Sprintf("%s (%s)", *sub.DisplayName, *sub.SubscriptionID))
		}
	}

	if len(items) == 0 {
		return fmt.Errorf("no subscriptions found")
	}

	// Let user select subscription
	selected, err := utils.SelectWithMenu(items, fmt.Sprintf("Select subscription for %s:", prompt))
	if err != nil {
		if strings.Contains(err.Error(), "cancelled by user") {
			fmt.Println("\nOperation cancelled by user")
			os.Exit(0)
		}
		return fmt.Errorf("failed to select subscription: %v", err)
	}

	// Extract subscription ID from selection
	subID, err = utils.ExtractIDFromParentheses(selected)
	if err != nil {
		return fmt.Errorf("failed to extract subscription ID: %v", err)
	}

	return nil
}

// SelectConnectionType prompts the user to select the type of connection.
func SelectConnectionType() (string, error) {
	return utils.SelectWithMenu([]string{"ssh", "tunnel"}, "Select connection type")
}

// GetAzureResources retrieves the necessary Azure resource configuration.
func GetAzureResources() (*ResourceConfig, error) {
	if err := ensureAuthenticated(); err != nil {
		return nil, fmt.Errorf("failed to authenticate: %v", err)
	}

	// Select subscription for Bastion host
	if err := selectSubscription("Bastion host"); err != nil {
		return nil, fmt.Errorf("failed to select subscription: %v", err)
	}

	// Get Bastion host details
	bastion, err := GetBastionDetails()
	if err != nil {
		return nil, fmt.Errorf("failed to get Bastion details: %v", err)
	}

	// Save Bastion subscription ID
	bastionSubID := subID

	// Select subscription for virtual machine
	if err := selectSubscription("virtual machine"); err != nil {
		return nil, fmt.Errorf("failed to select subscription: %v", err)
	}

	// Get target resource details
	target, err := GetTargetResource(ctx, subID, cred)
	if err != nil {
		return nil, fmt.Errorf("failed to get target resource: %v", err)
	}

	// Restore Bastion subscription ID for future use
	subID = bastionSubID

	return &ResourceConfig{
		BastionHost:    bastion,
		TargetResource: target,
	}, nil
}

// Connect establishes a connection to an Azure resource using the specified connection type.
func Connect(connectionType string, config *ResourceConfig) error {
	if config == nil {
		return fmt.Errorf("config cannot be nil")
	}

	switch connectionType {
	case "ssh":
		return connectSSH(config)
	case "tunnel":
		return connectTunnel(config)
	default:
		return fmt.Errorf("invalid connection type: %s", connectionType)
	}
}

func connectSSH(config *ResourceConfig) error {
	if config.Username == "" {
		username, err := utils.ReadInput("Enter username")
		if err != nil {
			return fmt.Errorf("failed to read username: %v", err)
		}
		config.Username = username
	}

	fmt.Println("Connecting via SSH...")
	// TODO: Implement SSH connection using Azure SDK
	// For now, use the existing CLI command
	args := []string{
		"network", "bastion", "ssh",
		"--name", config.BastionHost.Name,
		"--resource-group", config.BastionHost.ResourceGroup,
		"--target-resource-id", config.TargetResource.ID,
		"--auth-type", "password",
		"--username", config.Username,
	}

	if err := utils.AzureInteractiveCommand(args...); err != nil {
		return fmt.Errorf("failed to establish SSH connection: %v", err)
	}

	return nil
}

func connectTunnel(config *ResourceConfig) error {
	if config.RemotePort == 0 {
		remotePort, err := utils.ReadInput("Enter the resource port (default 22)")
		if err != nil {
			return fmt.Errorf("failed to read remote port: %v", err)
		}
		if remotePort == "" {
			remotePort = "22"
		}
		if _, err := fmt.Sscanf(remotePort, "%d", &config.RemotePort); err != nil {
			return fmt.Errorf("invalid remote port: %v", err)
		}
	}

	if config.LocalPort == 0 {
		localPort, err := utils.ReadInput("Enter the local port for the tunnel (e.g., 50022)")
		if err != nil {
			return fmt.Errorf("failed to read local port: %v", err)
		}
		if _, err := fmt.Sscanf(localPort, "%d", &config.LocalPort); err != nil {
			return fmt.Errorf("invalid local port: %v", err)
		}
	}

	fmt.Println("Establishing tunnel...")
	// TODO: Implement tunnel connection using Azure SDK
	// For now, use the existing CLI command
	args := []string{
		"network", "bastion", "tunnel",
		"--name", config.BastionHost.Name,
		"--resource-group", config.BastionHost.ResourceGroup,
		"--target-resource-id", config.TargetResource.ID,
		"--resource-port", fmt.Sprintf("%d", config.RemotePort),
		"--port", fmt.Sprintf("%d", config.LocalPort),
	}

	if err := utils.AzureInteractiveCommand(args...); err != nil {
		return fmt.Errorf("failed to establish tunnel: %v", err)
	}

	fmt.Printf("Tunnel established! You can now connect to localhost:%d\n", config.LocalPort)
	return nil
}
