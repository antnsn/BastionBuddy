package azure

import (
	"fmt"
	"time"

	"github.com/antnsn/BastionBuddy/internal/config"
	"github.com/antnsn/BastionBuddy/internal/tunnels"
)

// StartTunnel starts a new tunnel with the given configuration
func StartTunnel(resourceConfig *config.ResourceConfig, tunnelConfig *tunnels.Config) error {
	if err := ensureAuthenticated(); err != nil {
		return fmt.Errorf("failed to authenticate: %v", err)
	}

	// If no tunnel config is provided, create one from the resource config
	if tunnelConfig == nil {
		tunnelConfig = &tunnels.Config{
			Name:                  fmt.Sprintf("tunnel-%s", resourceConfig.TargetResource.Name),
			SubscriptionID:        resourceConfig.TargetResource.SubscriptionID,
			ResourceID:            resourceConfig.TargetResource.ID,
			ResourceName:          resourceConfig.TargetResource.Name,
			LocalPort:             resourceConfig.LocalPort,
			RemotePort:            resourceConfig.RemotePort,
			BastionName:           resourceConfig.BastionHost.Name,
			BastionResourceGroup:  resourceConfig.BastionHost.ResourceGroup,
			BastionSubscriptionID: resourceConfig.BastionHost.SubscriptionID,
			Username:              resourceConfig.Username,
			ConnectionType:        "tunnel",
			LastUsed:              time.Now(),
		}
	} else {
		// Update LastUsed time for existing config
		tunnelConfig.LastUsed = time.Now()
	}

	// Get the tunnel manager
	manager, err := GetTunnelManager()
	if err != nil {
		return fmt.Errorf("failed to get tunnel manager: %v", err)
	}

	// Save the configuration for future use
	if err := manager.configMgr.SaveConfig(*tunnelConfig); err != nil {
		return fmt.Errorf("failed to save tunnel configuration: %v", err)
	}

	// Start the tunnel
	_, err = manager.StartTunnel(
		tunnelConfig.SubscriptionID,
		tunnelConfig.ResourceID,
		tunnelConfig.ResourceName,
		tunnelConfig.LocalPort,
		tunnelConfig.RemotePort,
		tunnelConfig.BastionName,
		tunnelConfig.BastionResourceGroup,
		tunnelConfig.BastionSubscriptionID,
	)
	return err
}

// StartSavedTunnel starts a tunnel using a saved configuration
func StartSavedTunnel(tunnelName string) error {
	if err := ensureAuthenticated(); err != nil {
		return fmt.Errorf("failed to authenticate: %v", err)
	}

	manager, err := GetTunnelManager()
	if err != nil {
		return fmt.Errorf("failed to get tunnel manager: %v", err)
	}

	// Get all saved configurations
	configs := manager.GetSavedConfigs()

	// Find the requested tunnel configuration
	var tunnelConfig *tunnels.Config
	for _, config := range configs {
		if config.Name == tunnelName {
			tunnelConfig = &config
			break
		}
	}

	if tunnelConfig == nil {
		return fmt.Errorf("tunnel configuration '%s' not found", tunnelName)
	}

	// Start the tunnel with the saved configuration
	return StartTunnel(nil, tunnelConfig)
}

// StartSavedSSH starts an SSH connection using a saved configuration
func StartSavedSSH(configName string) error {
	if err := ensureAuthenticated(); err != nil {
		return fmt.Errorf("failed to authenticate: %v", err)
	}

	manager, err := GetTunnelManager()
	if err != nil {
		return fmt.Errorf("failed to get tunnel manager: %v", err)
	}

	// Get all saved configurations
	configs := manager.GetSavedConfigsByType("ssh")

	// Find the requested configuration
	var savedConfig *tunnels.Config
	for _, config := range configs {
		if config.Name == configName {
			savedConfig = &config
			break
		}
	}

	if savedConfig == nil {
		return fmt.Errorf("SSH configuration '%s' not found", configName)
	}

	// Create resource config from saved config
	resourceConfig := &config.ResourceConfig{
		BastionHost: &config.BastionHost{
			Name:           savedConfig.BastionName,
			ResourceGroup:  savedConfig.BastionResourceGroup,
			SubscriptionID: savedConfig.BastionSubscriptionID,
		},
		TargetResource: &config.TargetResource{
			ID:   savedConfig.ResourceID,
			Name: savedConfig.ResourceName,
		},
		Username:   savedConfig.Username,
		LocalPort:  savedConfig.LocalPort,
		RemotePort: savedConfig.RemotePort,
	}

	// Update the last used time
	savedConfig.LastUsed = time.Now()
	if err := manager.configMgr.SaveConfig(*savedConfig); err != nil {
		return fmt.Errorf("failed to update last used time: %v", err)
	}

	// Connect using the saved configuration and auth type
	return connectSSH(resourceConfig, savedConfig.AuthType)
}

// StartSavedRDP starts an RDP connection using a saved configuration
func StartSavedRDP(configName string) error {
	if err := ensureAuthenticated(); err != nil {
		return fmt.Errorf("failed to authenticate: %v", err)
	}

	manager, err := GetTunnelManager()
	if err != nil {
		return fmt.Errorf("failed to get tunnel manager: %v", err)
	}

	// Get all saved configurations
	configs := manager.GetSavedConfigs()

	// Find the requested configuration
	var savedConfig *tunnels.Config
	for _, config := range configs {
		if config.Name == configName {
			savedConfig = &config
			break
		}
	}

	if savedConfig == nil {
		return fmt.Errorf("configuration '%s' not found", configName)
	}

	// Create resource config from saved config
	resourceConfig := &config.ResourceConfig{
		BastionHost: &config.BastionHost{
			Name:           savedConfig.BastionName,
			ResourceGroup:  savedConfig.BastionResourceGroup,
			SubscriptionID: savedConfig.BastionSubscriptionID,
		},
		TargetResource: &config.TargetResource{
			ID:   savedConfig.ResourceID,
			Name: savedConfig.ResourceName,
		},
		Username:   savedConfig.Username,
		LocalPort:  savedConfig.LocalPort,
		RemotePort: savedConfig.RemotePort,
	}

	return connectRDP(resourceConfig)
}

// ListConfigurations lists saved configurations, optionally filtered by type
func ListConfigurations(connectionType string) error {
	manager, err := GetTunnelManager()
	if err != nil {
		return fmt.Errorf("failed to get tunnel manager: %v", err)
	}

	var configs []tunnels.Config
	if connectionType != "" {
		configs = manager.GetSavedConfigsByType(connectionType)
	} else {
		configs = manager.GetSavedConfigs()
	}

	if len(configs) == 0 {
		fmt.Println("No saved configurations found")
		return nil
	}

	// Group configurations by type
	tunnelConfigs := []tunnels.Config{}
	sshConfigs := []tunnels.Config{}
	rdpConfigs := []tunnels.Config{}

	for _, config := range configs {
		switch config.ConnectionType {
		case "ssh":
			sshConfigs = append(sshConfigs, config)
		case "rdp":
			rdpConfigs = append(rdpConfigs, config)
		default:
			tunnelConfigs = append(tunnelConfigs, config)
		}
	}

	// Print configurations by type
	if len(tunnelConfigs) > 0 && (connectionType == "" || connectionType == "tunnel") {
		fmt.Println("\nTunnel Configurations:")
		fmt.Println("---------------------")
		for _, config := range tunnelConfigs {
			fmt.Printf("Name: %s\n", config.Name)
			fmt.Printf("  Resource: %s\n", config.ResourceName)
			fmt.Printf("  Ports: local=%d, remote=%d\n", config.LocalPort, config.RemotePort)
			fmt.Printf("  Last Used: %s\n", config.LastUsed.Format("2006-01-02 15:04:05"))
			fmt.Println()
		}
	}

	if len(sshConfigs) > 0 && (connectionType == "" || connectionType == "ssh") {
		fmt.Println("\nSSH Configurations:")
		fmt.Println("-----------------")
		for _, config := range sshConfigs {
			fmt.Printf("Name: %s\n", config.Name)
			fmt.Printf("  Resource: %s\n", config.ResourceName)
			fmt.Printf("  Username: %s\n", config.Username)
			fmt.Printf("  Last Used: %s\n", config.LastUsed.Format("2006-01-02 15:04:05"))
			fmt.Println()
		}
	}

	if len(rdpConfigs) > 0 && (connectionType == "" || connectionType == "rdp") {
		fmt.Println("\nRDP Configurations:")
		fmt.Println("-----------------")
		for _, config := range rdpConfigs {
			fmt.Printf("Name: %s\n", config.Name)
			fmt.Printf("  Resource: %s\n", config.ResourceName)
			fmt.Printf("  Username: %s\n", config.Username)
			fmt.Printf("  Last Used: %s\n", config.LastUsed.Format("2006-01-02 15:04:05"))
			fmt.Println()
		}
	}

	return nil
}

// RunTunnelAction executes the specified tunnel action
func RunTunnelAction(_ *config.ResourceConfig, tunnelID string, action string) error {
	if err := ensureAuthenticated(); err != nil {
		return fmt.Errorf("failed to authenticate: %v", err)
	}

	// Execute the tunnel action
	manager, err := GetTunnelManager()
	if err != nil {
		return fmt.Errorf("failed to get tunnel manager: %v", err)
	}

	switch action {
	case "stop":
		return manager.StopTunnel(tunnelID)
	case "stop-all":
		return manager.StopAllTunnels()
	default:
		return fmt.Errorf("unknown action: %s", action)
	}
}
