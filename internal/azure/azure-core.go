package azure

import (
	"context"
	"fmt"
	"time"

	"runtime"

	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/subscription/armsubscription"
	"github.com/antnsn/BastionBuddy/internal/config"
	"github.com/antnsn/BastionBuddy/internal/tunnels"
	"github.com/antnsn/BastionBuddy/internal/utils"
)

// init initializes the Azure package
func init() {
	// Initialize other things if needed, but don't check auth here
}

// ensureAuthenticated ensures the user is logged in to Azure
func ensureAuthenticated() error {
	// First check if we're logged into Azure CLI
	isLoggedIn, err := utils.CheckAzureLogin()
	if err != nil {
		return fmt.Errorf("failed to check Azure login status: %v", err)
	}

	// If not logged in, start the login process
	if !isLoggedIn {
		if err := utils.LoginToAzure(); err != nil {
			return fmt.Errorf("failed to login to Azure: %v", err)
		}
	}

	// Now check Azure SDK credentials
	cred, err := GetAzureCredential()
	if err != nil {
		return fmt.Errorf("authentication failed: %v", err)
	}
	if cred == nil {
		return fmt.Errorf("no valid Azure credentials found")
	}
	return nil
}

// CheckDependencies checks if all required dependencies are installed
func CheckDependencies() error {
	return utils.CheckDependencies()
}

// Cleanup performs any necessary cleanup operations.
func Cleanup() error {
	return nil
}

// getSubscriptions retrieves available Azure subscriptions using the Azure SDK
func getSubscriptions(ctx context.Context, cred *azidentity.DefaultAzureCredential) ([]*armsubscription.Subscription, error) {
	client, err := armsubscription.NewSubscriptionsClient(cred, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create subscriptions client: %v", err)
	}

	pager := client.NewListPager(nil)
	var subscriptions []*armsubscription.Subscription

	for pager.More() {
		page, err := pager.NextPage(ctx)
		if err != nil {
			return nil, fmt.Errorf("failed to get next page: %v", err)
		}

		for _, sub := range page.Value {
			if sub.DisplayName != nil && sub.SubscriptionID != nil {
				subscriptions = append(subscriptions, sub)
			}
		}
	}

	if len(subscriptions) == 0 {
		return nil, fmt.Errorf("no subscriptions found")
	}

	return subscriptions, nil
}

// SelectConnectionType prompts the user to select the type of connection.
func SelectConnectionType() (ConnectionType, error) {
	var items []string
	items = append(items, string(SSH))
	if runtime.GOOS == "windows" {
		items = append(items, string(RDP))
	}
	items = append(items, string(Tunnel))

	selected, err := utils.SelectWithMenu(items, "Select connection type")
	if err != nil {
		return "", fmt.Errorf("failed to select connection type: %v", err)
	}
	return ConnectionType(selected), nil
}

// GetAzureResources retrieves the necessary Azure resource configuration.
func GetAzureResources() (*config.ResourceConfig, error) {
	if err := ensureAuthenticated(); err != nil {
		return nil, err
	}

	// Get Azure credentials
	cred, err := GetAzureCredential()
	if err != nil {
		return nil, fmt.Errorf("failed to get Azure credentials: %v", err)
	}

	// Get subscription ID
	ctx := context.Background()
	subID, err := getSubscriptionID(ctx, cred, "Select Azure subscription for Bastion host")
	if err != nil {
		return nil, fmt.Errorf("failed to get subscription ID: %v", err)
	}

	// Get Bastion host details
	bastionHost, err := GetBastionDetails(ctx, cred, subID)
	if err != nil {
		return nil, fmt.Errorf("failed to get Bastion host details: %v", err)
	}

	// Get target resource details
	targetResource, err := GetTargetResource(ctx, cred, subID)
	if err != nil {
		return nil, fmt.Errorf("failed to get target resource details: %v", err)
	}

	// Get username
	username, err := utils.ReadInput("Enter username")
	if err != nil {
		return nil, fmt.Errorf("failed to read username: %v", err)
	}

	// Get local port
	localPort := 2222 // Default local port

	// Get remote port
	remotePort := 22 // Default remote port for SSH

	return &config.ResourceConfig{
		Username:       username,
		LocalPort:      localPort,
		RemotePort:     remotePort,
		BastionHost:    bastionHost,
		TargetResource: targetResource,
	}, nil
}

// InitiateConnection handles the complete connection flow
func InitiateConnection() error {
	for {
		// Step 1: User has already selected Connect to get here

		// Step 2: Get connection type (SSH/RDP/Tunnel)
		connectionType, err := SelectConnectionType()
		if err == utils.ErrReturnToMain {
			return nil // Return to main menu
		}
		if err != nil {
			return fmt.Errorf("failed to select connection type: %v", err)
		}

		// Step 3: Select subscription for Bastion host
		ctx := context.Background()
		cred, err := azidentity.NewDefaultAzureCredential(nil)
		if err != nil {
			return fmt.Errorf("failed to create credentials: %v", err)
		}

		bastionSubscriptionID, err := getSubscriptionID(ctx, cred, "Select Azure subscription for Bastion host")
		if err == utils.ErrReturnToMain {
			return nil // Return to main menu
		}
		if err != nil {
			return fmt.Errorf("failed to select Bastion subscription: %v", err)
		}

		// Step 4: Get Bastion host details
		bastionHost, err := GetBastionDetails(ctx, cred, bastionSubscriptionID)
		if err != nil {
			return fmt.Errorf("failed to get Bastion details: %v", err)
		}

		// Step 5: Select target machine subscription
		targetSubscriptionID, err := getSubscriptionID(ctx, cred, "Select Azure subscription for target resource")
		if err == utils.ErrReturnToMain {
			return nil // Return to main menu
		}
		if err != nil {
			return fmt.Errorf("failed to select target subscription: %v", err)
		}

		// Step 6: Get target resource details
		targetResource, err := GetTargetResource(ctx, cred, targetSubscriptionID)
		if err != nil {
			return fmt.Errorf("failed to get target resource: %v", err)
		}

		// Step 7: Get port configurations based on connection type
		config := &config.ResourceConfig{
			BastionHost:    bastionHost,
			TargetResource: targetResource,
		}

		// Get username for SSH connections
		if connectionType == SSH {
			username, err := utils.ReadInput("Enter username for SSH connection")
			if err != nil {
				return fmt.Errorf("failed to get username: %v", err)
			}
			if username == "" {
				return fmt.Errorf("username is required for SSH connections")
			}
			config.Username = username
		}

		if connectionType == Tunnel {
			var defaultRemotePort int
			switch connectionType {
			case SSH:
				defaultRemotePort = 22
			case RDP:
				defaultRemotePort = 3389
			default:
				defaultRemotePort = 0
			}

			// Step 7.1.1: Get target port
			portPrompt := "Enter target resource port (e.g., 22 for SSH, 3389 for RDP, 80 for HTTP, 443 for HTTPS)"
			if defaultRemotePort > 0 {
				portPrompt = fmt.Sprintf("Enter target resource port (default: %d)", defaultRemotePort)
			}
			remotePort, err := utils.GetUserInputInt(portPrompt)
			if err != nil {
				return fmt.Errorf("failed to get remote port: %v", err)
			}
			if remotePort == 0 && defaultRemotePort > 0 {
				remotePort = defaultRemotePort
			}
			config.RemotePort = remotePort

			// Step 7.1.2: Get local port
			localPortPrompt := fmt.Sprintf("Enter local port (e.g., %d to match remote port, or any available local port)", remotePort)
			localPort, err := utils.GetUserInputInt(localPortPrompt)
			if err != nil {
				return fmt.Errorf("failed to get local port: %v", err)
			}
			config.LocalPort = localPort

			// Create and start the tunnel
			tunnelConfig := &tunnels.Config{
				Name:                  fmt.Sprintf("tunnel-%s", config.TargetResource.Name),
				SubscriptionID:        config.TargetResource.SubscriptionID,
				ResourceID:            config.TargetResource.ID,
				ResourceName:          config.TargetResource.Name,
				LocalPort:             config.LocalPort,
				RemotePort:            config.RemotePort,
				Command:               "",
				Args:                  nil,
				LastUsed:              time.Now(),
				BastionName:           config.BastionHost.Name,
				BastionResourceGroup:  config.BastionHost.ResourceGroup,
				BastionSubscriptionID: config.BastionHost.SubscriptionID,
				ConnectionType:        "tunnel",
				Username:              config.Username,
			}

			if _, err := StartTunnel(config, tunnelConfig); err != nil {
				return fmt.Errorf("failed to start tunnel: %v", err)
			}

			fmt.Printf("\nTunnel created successfully! Local port %d is now forwarding to remote port %d\n", localPort, remotePort)
			fmt.Printf("Use 'manage-tunnels' from the main menu to view and manage active tunnels\n")
			return nil
		}

		// Handle SSH/RDP connection
		if err := establishConnection(connectionType, config); err != nil {
			return fmt.Errorf("failed to establish %s connection: %v", connectionType, err)
		}

	}
}

// establishConnection establishes a connection to an Azure resource using the specified connection type.
func establishConnection(connectionType ConnectionType, config *config.ResourceConfig) error {
	if err := ensureAuthenticated(); err != nil {
		return err
	}

	if connectionType == RDP && runtime.GOOS != "windows" {
		return fmt.Errorf("RDP connections are only supported on Windows")
	}

	if config == nil {
		var err error
		config, err = GetAzureResources()
		if err != nil {
			return fmt.Errorf("failed to get Azure resources: %v", err)
		}
	}

	switch connectionType {
	case SSH:
		if err := connectSSH(config, ""); err != nil {
			return err
		}
	case Tunnel:
		tunnelConfig := &tunnels.Config{
			Name:                  fmt.Sprintf("tunnel-%s", config.TargetResource.Name),
			SubscriptionID:        config.TargetResource.SubscriptionID,
			ResourceID:            config.TargetResource.ID,
			ResourceName:          config.TargetResource.Name,
			LocalPort:             config.LocalPort,
			RemotePort:            config.RemotePort,
			Command:               "",
			Args:                  nil,
			LastUsed:              time.Now(),
			BastionName:           config.BastionHost.Name,
			BastionResourceGroup:  config.BastionHost.ResourceGroup,
			BastionSubscriptionID: config.BastionHost.SubscriptionID,
			ConnectionType:        "tunnel",
			Username:              config.Username,
		}
		if _, err := StartTunnel(config, tunnelConfig); err != nil {
			return err
		}
	case RDP:
		if err := connectRDP(config, nil); err != nil {
			return err
		}
	default:
		return fmt.Errorf("invalid connection type: %s", connectionType)
	}

	return nil
}

func connectSSH(config *config.ResourceConfig, savedAuthType string) error {
	if config == nil {
		return fmt.Errorf("no configuration provided")
	}

	if config.Username == "" {
		return fmt.Errorf("username is required")
	}

	var authType string

	if savedAuthType != "" {
		// Use the saved auth type
		authType = savedAuthType
	} else {
		// Let user select auth type for new connections
		authType, _ = utils.SelectWithMenu([]string{"AAD", "password"}, "Select authentication type")
	}

	// Save the SSH configuration only if it's a new connection
	if savedAuthType == "" {
		manager, err := GetTunnelManager()
		if err != nil {
			return fmt.Errorf("failed to get tunnel manager: %v", err)
		}

		sshConfig := tunnels.Config{
			Name:                  fmt.Sprintf("ssh-%s", config.TargetResource.Name),
			SubscriptionID:        config.TargetResource.SubscriptionID,
			ResourceID:            config.TargetResource.ID,
			ResourceName:          config.TargetResource.Name,
			BastionName:           config.BastionHost.Name,
			BastionResourceGroup:  config.BastionHost.ResourceGroup,
			BastionSubscriptionID: config.BastionHost.SubscriptionID,
			Username:              config.Username,
			ConnectionType:        "ssh",
			LastUsed:              time.Now(),
			AuthType:              authType,
		}

		if err := manager.configMgr.SaveConfig(sshConfig); err != nil {
			return fmt.Errorf("failed to save SSH configuration: %v", err)
		}
	}

	args := []string{
		"network", "bastion", "ssh",
		"--subscription", config.BastionHost.SubscriptionID,
		"--resource-group", config.BastionHost.ResourceGroup,
		"--name", config.BastionHost.Name,
		"--target-resource-id", config.TargetResource.ID,
		"--auth-type", authType,
		"--username", config.Username,
	}

	return utils.AzureInteractiveCommand(args...)
}

func connectRDP(config *config.ResourceConfig, savedConfig *tunnels.Config) error {
	if runtime.GOOS != "windows" {
		return fmt.Errorf("RDP connections are only supported on Windows")
	}

	if config == nil {
		return fmt.Errorf("no configuration provided")
	}

	if config.Username == "" {
		return fmt.Errorf("username is required")
	}

	var enableMFA bool

	if savedConfig != nil {
		// Use saved MFA setting
		enableMFA = savedConfig.EnableMFA
	} else {
		// Ask user if they want to enable MFA
		mfaChoice, _ := utils.SelectWithMenu([]string{"No", "Yes"}, "Enable Multi-Factor Authentication (MFA)?")
		enableMFA = mfaChoice == "Yes"

		// Save the RDP configuration for new connections
		manager, err := GetTunnelManager()
		if err != nil {
			return fmt.Errorf("failed to get tunnel manager: %v", err)
		}

		rdpConfig := tunnels.Config{
			Name:                  fmt.Sprintf("rdp-%s", config.TargetResource.Name),
			SubscriptionID:        config.TargetResource.SubscriptionID,
			ResourceID:            config.TargetResource.ID,
			ResourceName:          config.TargetResource.Name,
			BastionName:           config.BastionHost.Name,
			BastionResourceGroup:  config.BastionHost.ResourceGroup,
			BastionSubscriptionID: config.BastionHost.SubscriptionID,
			Username:              config.Username,
			ConnectionType:        "rdp",
			LastUsed:              time.Now(),
			EnableMFA:             enableMFA,
		}

		if err := manager.configMgr.SaveConfig(rdpConfig); err != nil {
			return fmt.Errorf("failed to save RDP configuration: %v", err)
		}
	}

	args := []string{
		"network", "bastion", "rdp",
		"--subscription", config.BastionHost.SubscriptionID,
		"--resource-group", config.BastionHost.ResourceGroup,
		"--name", config.BastionHost.Name,
		"--target-resource-id", config.TargetResource.ID,
		"--auth-type", "AAD",
		"--username", config.Username,
	}

	if enableMFA {
		args = append(args, "--enable-mfa")
	}

	return utils.AzureInteractiveCommand(args...)
}

// SelectInitialAction prompts the user to select the initial action
func SelectInitialAction() (string, error) {
	items := []string{"Create new connection"}

	// Only show manage-tunnels if there are active tunnels
	if manager, err := GetTunnelManager(); err == nil {
		if len(manager.ListTunnels()) > 0 {
			items = append(items, "Manage active tunnels")
		}
	}

	items = append(items, "Exit BastionBuddy")

	action, err := utils.SelectWithMenu(items, "What would you like to do?")
	if err != nil {
		return "", fmt.Errorf("failed to select action: %v", err)
	}

	// Map friendly names back to internal action names
	switch action {
	case "Create new connection":
		return "connect", nil
	case "Manage active tunnels":
		return "manage-tunnels", nil
	case "Exit BastionBuddy":
		return "exit", nil
	default:
		return action, nil
	}
}

// handleInitialAction handles the selected initial action
func handleInitialAction(action string, _ *config.ResourceConfig) error {
	for {
		switch action {
		case "connect":
			err := InitiateConnection()
			if err == utils.ErrReturnToMain {
				// Get a new action selection
				newAction, err := SelectInitialAction()
				if err != nil {
					return err
				}
				action = newAction
				continue
			}
			return err
		case "manage-tunnels":
			err := manageTunnels()
			if err == utils.ErrReturnToMain {
				// Get a new action selection
				newAction, err := SelectInitialAction()
				if err != nil {
					return err
				}
				action = newAction
				continue
			}
			return err
		case "exit":
			return nil
		default:
			return fmt.Errorf("unknown action: %s", action)
		}
	}
}

func manageTunnels() error {
	manager, err := GetTunnelManager()
	if err != nil {
		return fmt.Errorf("failed to get tunnel manager: %v", err)
	}

	tunnels := manager.ListTunnels()

	// Create menu items for each active tunnel
	var items []string
	tunnelMap := make(map[string]*TunnelInfo)
	for _, t := range tunnels {
		item := fmt.Sprintf("Terminate: %s (Local:%d → Remote:%d) - %s [Active: %s]",
			t.ID[:8], t.LocalPort, t.RemotePort, t.ResourceName,
			time.Since(t.StartTime).Round(time.Second))
		items = append(items, item)
		tunnelMap[item] = t
	}
	items = append(items, "Return to main menu")

	selected, err := utils.SelectWithMenu(items, "Select a tunnel to terminate, or return to main menu")
	if err != nil {
		if err == utils.ErrReturnToMain {
			return nil
		}
		return fmt.Errorf("failed to select tunnel: %v", err)
	}

	if selected == "Return to main menu" {
		return nil
	}

	if tunnel, ok := tunnelMap[selected]; ok {
		if err := manager.StopTunnel(tunnel.ID); err != nil {
			return fmt.Errorf("failed to stop tunnel: %v", err)
		}
		return nil
	}

	return fmt.Errorf("invalid selection")
}

// InitiateAction is the exported function that handles the initial action
func InitiateAction(action string, config *config.ResourceConfig) error {
	return handleInitialAction(action, config)
}

// getSubscriptionID retrieves the subscription ID from the Azure CLI.
func getSubscriptionID(ctx context.Context, cred *azidentity.DefaultAzureCredential, prompt string) (string, error) {
	// Get the list of subscriptions
	subs, err := getSubscriptions(ctx, cred)
	if err != nil {
		return "", err
	}

	// If there's only one subscription, use it
	if len(subs) == 1 {
		return *subs[0].SubscriptionID, nil
	}

	// Create menu items
	var items []string
	subMap := make(map[string]*armsubscription.Subscription)
	for _, sub := range subs {
		item := fmt.Sprintf("%s | ID: %s | State: %s",
			*sub.DisplayName,
			*sub.SubscriptionID,
			*sub.State)
		items = append(items, item)
		subMap[item] = sub
	}

	// Let the user select a subscription
	selected, err := utils.SelectWithMenu(items, prompt+" (type to filter)")
	if err != nil {
		return "", fmt.Errorf("failed to select subscription: %v", err)
	}

	// Return the selected subscription ID
	return *subMap[selected].SubscriptionID, nil
}
