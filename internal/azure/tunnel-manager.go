package azure

import (
	"fmt"
	"os/exec"
	"time"

	"github.com/antnsn/BastionBuddy/internal/tunnels"
	"github.com/antnsn/BastionBuddy/internal/utils"
	"github.com/google/uuid"
)

// TunnelInfo contains information about a tunnel connection
type TunnelInfo struct {
	ID                    string
	LocalPort             int
	RemotePort            int
	ResourceID            string
	ResourceName          string
	SubscriptionID        string
	BastionName           string
	BastionResourceGroup  string
	BastionSubscriptionID string
	StartTime             time.Time
	Status                string
	cmd                   *exec.Cmd
}

// TunnelManager manages tunnel connections
type TunnelManager struct {
	tunnels   map[string]*TunnelInfo
	configMgr *tunnels.Manager
}

// ListTunnels returns a list of all active tunnels
func (tm *TunnelManager) ListTunnels() []*TunnelInfo {
	tunnels := make([]*TunnelInfo, 0, len(tm.tunnels))
	for _, tunnel := range tm.tunnels {
		tunnels = append(tunnels, tunnel)
	}
	return tunnels
}

// GetSavedConfigs returns a list of saved tunnel configurations
func (tm *TunnelManager) GetSavedConfigs() []tunnels.Config {
	return tm.configMgr.GetSavedConfigs()
}

// GetSavedConfigsByType returns saved configurations of a specific type
func (tm *TunnelManager) GetSavedConfigsByType(connectionType string) []tunnels.Config {
	return tm.configMgr.GetSavedConfigsByType(connectionType)
}

// StopTunnel stops a specific tunnel
func (tm *TunnelManager) StopTunnel(id string) error {
	tunnel, exists := tm.tunnels[id]
	if !exists {
		return fmt.Errorf("tunnel %s not found", id)
	}

	if err := tm.stopTunnelProcess(tunnel); err != nil {
		return fmt.Errorf("failed to stop tunnel process: %v", err)
	}

	// Remove from in-memory state
	delete(tm.tunnels, id)

	// Remove from persistent storage
	if err := tm.configMgr.RemoveActive(id); err != nil {
		return fmt.Errorf("failed to remove tunnel from storage: %v", err)
	}

	return nil
}

// StopAllTunnels stops all active tunnels
func (tm *TunnelManager) StopAllTunnels() error {
	var lastErr error
	for id, tunnel := range tm.tunnels {
		if err := tm.stopTunnelProcess(tunnel); err != nil {
			lastErr = fmt.Errorf("failed to stop tunnel %s: %v", id, err)
			fmt.Printf("Warning: %v\n", lastErr)
			continue
		}

		// Remove from persistent storage
		if err := tm.configMgr.RemoveActive(id); err != nil {
			lastErr = fmt.Errorf("failed to remove tunnel %s from storage: %v", id, err)
			fmt.Printf("Warning: %v\n", lastErr)
		}
	}

	// Clear in-memory state
	tm.tunnels = make(map[string]*TunnelInfo)

	return lastErr
}

// stopTunnelProcess stops the process associated with a tunnel
func (tm *TunnelManager) stopTunnelProcess(tunnel *TunnelInfo) error {
	if tunnel.cmd != nil && tunnel.cmd.Process != nil {
		if err := tunnel.cmd.Process.Kill(); err != nil {
			return fmt.Errorf("failed to kill tunnel process: %v", err)
		}
	}
	return nil
}

// StartTunnel starts a new tunnel connection
func (tm *TunnelManager) StartTunnel(subscriptionID string, resourceID string, resourceName string, localPort int, remotePort int, bastionName string, bastionResourceGroup string, bastionSubscriptionID string) (*TunnelInfo, error) {
	// Create a new tunnel info
	tunnel := &TunnelInfo{
		ID:                    uuid.New().String(),
		LocalPort:             localPort,
		RemotePort:            remotePort,
		ResourceID:            resourceID,
		ResourceName:          resourceName,
		SubscriptionID:        subscriptionID,
		BastionName:           bastionName,
		BastionResourceGroup:  bastionResourceGroup,
		BastionSubscriptionID: bastionSubscriptionID,
		StartTime:             time.Now(),
		Status:                "starting",
	}

	// Prepare the Azure command
	cmd := utils.PrepareAzureCommand("network", "bastion", "tunnel",
		"--subscription", subscriptionID,
		"--target-resource-id", resourceID,
		"--resource-port", fmt.Sprintf("%d", remotePort),
		"--port", fmt.Sprintf("%d", localPort))

	cmd.SysProcAttr = utils.GetSysProcAttr()
	tunnel.cmd = cmd

	// Start the tunnel process
	if err := cmd.Start(); err != nil {
		return nil, fmt.Errorf("failed to start tunnel: %v", err)
	}

	// Save the tunnel info
	tm.tunnels[tunnel.ID] = tunnel

	// Save the tunnel configuration
	activeTunnel := &tunnels.Active{
		ID:                    tunnel.ID,
		LocalPort:             localPort,
		RemotePort:            remotePort,
		ResourceID:            resourceID,
		ResourceName:          resourceName,
		SubscriptionID:        subscriptionID,
		BastionName:           bastionName,
		BastionResourceGroup:  bastionResourceGroup,
		BastionSubscriptionID: bastionSubscriptionID,
		StartTime:             time.Now(),
		Status:                "active",
	}
	if err := tm.configMgr.SaveActive(*activeTunnel); err != nil {
		// Clean up if saving fails
		_ = tm.stopTunnelProcess(tunnel)
		delete(tm.tunnels, tunnel.ID)
		return nil, fmt.Errorf("failed to save tunnel configuration: %v", err)
	}

	// Save the tunnel config for future use
	tunnelConfig := &tunnels.Config{
		Name:                  fmt.Sprintf("tunnel-%s", resourceName),
		SubscriptionID:        subscriptionID,
		ResourceID:            resourceID,
		ResourceName:          resourceName,
		LocalPort:             localPort,
		RemotePort:            remotePort,
		BastionName:           bastionName,
		BastionResourceGroup:  bastionResourceGroup,
		BastionSubscriptionID: bastionSubscriptionID,
		LastUsed:              time.Now(),
	}
	if err := tm.configMgr.SaveConfig(*tunnelConfig); err != nil {
		// Log the error but don't fail the tunnel creation
		fmt.Printf("Warning: failed to save tunnel configuration: %v\n", err)
	}

	return tunnel, nil
}

// SaveActive saves an active tunnel configuration
func (tm *TunnelManager) SaveActive(tunnel tunnels.Active) error {
	return tm.configMgr.SaveActive(tunnel)
}
