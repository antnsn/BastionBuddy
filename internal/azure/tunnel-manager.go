package azure

import (
	"bytes"
	"fmt"
	"net"
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
	PID                   int
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
	fmt.Printf("Stopping tunnel process for ID: %s\n", tunnel.ID)
	if tunnel.cmd != nil && tunnel.cmd.Process != nil {
		if err := tunnel.cmd.Process.Kill(); err != nil {
			fmt.Printf("Failed to kill tunnel process for ID %s: %v\n", tunnel.ID, err)
			return fmt.Errorf("failed to kill tunnel process: %v", err)
		}
		fmt.Printf("Successfully killed tunnel process for ID: %s\n", tunnel.ID)
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
		"--subscription", bastionSubscriptionID, // Use bastion's subscription ID
		"--target-resource-id", resourceID,
		"--resource-port", fmt.Sprintf("%d", remotePort),
		"--port", fmt.Sprintf("%d", localPort),
		"--name", bastionName,
		"--resource-group", bastionResourceGroup)

	// Print command for debugging
	// fmt.Printf("Starting tunnel with command: %s %v\n", cmd.Path, cmd.Args)

	// Capture output for debugging
	var outputBuffer bytes.Buffer
	cmd.Stdout = &outputBuffer
	cmd.Stderr = &outputBuffer

	cmd.SysProcAttr = utils.GetSysProcAttr()
	tunnel.cmd = cmd

	// Start the tunnel process
	if err := cmd.Start(); err != nil {
		return nil, fmt.Errorf("failed to start tunnel: %v", err)
	}

	// Store the PID
	tunnel.PID = cmd.Process.Pid

	// Save the tunnel info
	tm.tunnels[tunnel.ID] = tunnel

	// Wait a moment and check if the process is still running
	time.Sleep(2 * time.Second)
	if cmd.Process == nil || cmd.ProcessState != nil && cmd.ProcessState.Exited() {
		tunnel.Status = "failed"
		output := outputBuffer.String()
		return nil, fmt.Errorf("tunnel process failed to start or exited immediately: %s", output)
	}

	// Check if the port is actually being listened on
	conn, err := net.DialTimeout("tcp", fmt.Sprintf("127.0.0.1:%d", localPort), time.Second)
	if err != nil {
		tunnel.Status = "failed"
		output := outputBuffer.String()
		if err := cmd.Process.Kill(); err != nil {
			fmt.Printf("Warning: failed to kill tunnel process: %v\n", err)
		}
		return nil, fmt.Errorf("tunnel port %d is not listening after startup: %v\nOutput: %s", localPort, err, output)
	}

	// Close the connection and check for errors
	if err := conn.Close(); err != nil {
		return nil, fmt.Errorf("failed to close connection: %v", err)
	}

	// Update status to running
	tunnel.Status = "running"

	// Print connection command
	tm.PrintConnectionCommand(tunnel)

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
		PID:                   tunnel.PID, // Include PID
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

// PrintConnectionCommand prints the command to connect to a tunnel
func (tm *TunnelManager) PrintConnectionCommand(tunnel *TunnelInfo) {
	// For localhost tunnels, we want to skip host key checking since the key will change
	// with different local ports
	fmt.Printf("\nTunnel activated:\n")
	fmt.Printf("Connection available at: localhost\n")
	fmt.Printf("Port: %d\n", tunnel.LocalPort)
	fmt.Printf("\n")
}

// SaveActive saves an active tunnel configuration
func (tm *TunnelManager) SaveActive(tunnel tunnels.Active) error {
	return tm.configMgr.SaveActive(tunnel)
}
