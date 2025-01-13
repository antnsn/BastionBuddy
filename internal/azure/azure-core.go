package azure

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/antnsn/BastionBuddy/internal/cache"
	"github.com/antnsn/BastionBuddy/internal/utils"
	"github.com/antnsn/BastionBuddy/internal/welcome"
)

var cacheInstance *cache.Cache

func init() {
	cacheDir := filepath.Join(os.TempDir(), "bastionbuddy-cache")
	cacheInstance = cache.NewCache(cacheDir, 1*time.Hour)
}

// CheckDependencies verifies that all required Azure CLI dependencies are installed
func CheckDependencies() error {
	return utils.CheckDependencies()
}

// Cleanup performs any necessary cleanup operations
func Cleanup() {
	if cacheInstance != nil {
		_ = cacheInstance.Cleanup()
	}
}

// ensureAuthenticated ensures the user is authenticated with Azure CLI
func ensureAuthenticated() error {
	if _, err := utils.AzureCommand("account", "show", "--query", "id", "-o", "tsv"); err != nil {
		fmt.Println("Authenticating with Azure...")
		if _, err := utils.AzureCommand("login", "--use-device-code"); err != nil {
			return fmt.Errorf("authentication failed: %v", err)
		}
	}
	return nil
}

// selectSubscription prompts the user to select an Azure subscription
func selectSubscription(prompt string) (string, error) {
	fmt.Printf("Fetching subscriptions for %s...\n", prompt)

	output, err := utils.AzureCommand("account", "list", "--query", "[].{Name:name, Id:id}", "-o", "json")
	if err != nil {
		return "", fmt.Errorf("failed to list subscriptions: %v", err)
	}

	var subs []struct {
		Name string `json:"Name"`
		ID   string `json:"Id"`
	}
	if err := json.Unmarshal(output, &subs); err != nil {
		return "", fmt.Errorf("failed to parse subscriptions: %v", err)
	}

	var options []string
	for _, sub := range subs {
		options = append(options, fmt.Sprintf("%s (%s)", sub.Name, sub.ID))
	}

	selected, err := utils.SelectWithMenu(options, fmt.Sprintf("Select subscription for %s:", prompt))
	if err != nil {
		return "", fmt.Errorf("failed to select subscription: %v", err)
	}

	return utils.ExtractIDFromParentheses(selected)
}

// SelectConnectionType prompts the user to select a connection type
func SelectConnectionType() (string, error) {
	welcome.ShowWelcome()
	return utils.SelectWithMenu([]string{"ssh", "tunnel"}, "Select connection type:")
}

// GetAzureResources gathers all necessary Azure resource information
func GetAzureResources() (*AzureConfig, error) {
	if err := ensureAuthenticated(); err != nil {
		return nil, err
	}

	config := &AzureConfig{}

	// Select Bastion subscription
	bastionSubID, err := selectSubscription("Bastion host")
	if err != nil {
		return nil, fmt.Errorf("failed to select Bastion subscription: %v", err)
	}

	if err := utils.AzureSetSubscription(bastionSubID); err != nil {
		return nil, fmt.Errorf("failed to set Bastion subscription: %v", err)
	}

	// Get Bastion details
	bastionName, bastionRG, err := GetBastionDetails(bastionSubID)
	if err != nil {
		return nil, fmt.Errorf("failed to get bastion details: %v", err)
	}
	config.BastionName = bastionName
	config.BastionResourceGroup = bastionRG

	// Get target resource ID
	resourceID, err := GetTargetResource()
	if err != nil {
		return nil, fmt.Errorf("failed to get target resource: %v", err)
	}
	config.ResourceID = resourceID

	// Switch back to Bastion subscription for the connection
	if err := utils.AzureSetSubscription(bastionSubID); err != nil {
		return nil, fmt.Errorf("failed to set subscription for connection: %v", err)
	}

	return config, nil
}

// Connect establishes a connection using the specified type and configuration
func Connect(connectionType string, config *AzureConfig) error {
	switch connectionType {
	case "ssh":
		return connectSSH(config)
	case "tunnel":
		return connectTunnel(config)
	default:
		return fmt.Errorf("unsupported connection type: %s", connectionType)
	}
}

func connectSSH(config *AzureConfig) error {
	if config.Username == "" {
		username, err := utils.ReadInput("Enter username: ")
		if err != nil {
			return fmt.Errorf("failed to read username: %v", err)
		}
		if username == "" {
			return fmt.Errorf("username cannot be empty")
		}
		config.Username = username
	}

	fmt.Printf("Debug: Connecting to resource '%s' via SSH with username '%s'\n", config.ResourceID, config.Username)
	fmt.Printf("Debug: Using Bastion host '%s' in resource group '%s'\n", config.BastionName, config.BastionResourceGroup)
	
	err := utils.AzureInteractiveCommand("network", "bastion", "ssh",
		"--name", config.BastionName,
		"--resource-group", config.BastionResourceGroup,
		"--target-resource-id", config.ResourceID,
		"--auth-type", "password",
		"--username", config.Username)
	if err != nil {
		return fmt.Errorf("failed to establish SSH connection: %v\nPlease check:\n"+
			"1. The Bastion host name and resource group are correct\n"+
			"2. You have the necessary permissions on both the Bastion and target resources\n"+
			"3. The target resource is running and accessible\n"+
			"4. The username is correct and exists on the target resource", err)
	}
	return nil
}

func connectTunnel(config *AzureConfig) error {
	if config.LocalPort == 0 {
		input, err := utils.ReadInput("Enter the local port for the tunnel (e.g., 50022): ")
		if err != nil {
			return fmt.Errorf("failed to read local port: %v", err)
		}
		var port int
		if _, err := fmt.Sscanf(input, "%d", &port); err != nil {
			return fmt.Errorf("invalid local port: %v", err)
		}
		config.LocalPort = port
	}

	if config.RemotePort == 0 {
		input, err := utils.ReadInput("Enter the resource port (default 22): ")
		if err != nil {
			return fmt.Errorf("failed to read remote port: %v", err)
		}
		if input == "" {
			config.RemotePort = 22
		} else {
			var port int
			if _, err := fmt.Sscanf(input, "%d", &port); err != nil {
				return fmt.Errorf("invalid remote port: %v", err)
			}
			config.RemotePort = port
		}
	}

	fmt.Printf("Creating tunnel from localhost:%d to %s:%d...\n",
		config.LocalPort, config.ResourceID, config.RemotePort)

	_, err := utils.AzureCommand("network", "bastion", "tunnel",
		"--name", config.BastionName,
		"--resource-group", config.BastionResourceGroup,
		"--target-resource-id", config.ResourceID,
		"--resource-port", fmt.Sprintf("%d", config.RemotePort),
		"--port", fmt.Sprintf("%d", config.LocalPort))
	if err != nil {
		return fmt.Errorf("failed to establish tunnel: %v\nPlease check:\n"+
			"1. The Bastion host name and resource group are correct\n"+
			"2. You have the necessary permissions on both the Bastion and target resources\n"+
			"3. The target resource is running and accessible", err)
	}

	fmt.Printf("Tunnel established! You can now connect to localhost:%d\n", config.LocalPort)
	return nil
}
