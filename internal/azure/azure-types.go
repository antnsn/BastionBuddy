// Package azure provides functionality for interacting with Azure resources
// through the Azure CLI, specifically for Bastion connections.
package azure

import (
	"fmt"
	"github.com/antnsn/BastionBuddy/internal/config"
)

var debugMode bool

// SetDebugMode enables or disables debug output
func SetDebugMode(enabled bool) {
	debugMode = enabled
}

// debugPrintf prints debug messages if debug mode is enabled
func debugPrintf(format string, args ...interface{}) {
	if debugMode {
		fmt.Printf(format, args...)
	}
}

// ConnectionType represents the type of connection to establish.
type ConnectionType string

const (
	// SSH represents an SSH connection.
	SSH ConnectionType = "ssh"
	// RDP represents an RDP connection.
	RDP ConnectionType = "rdp"
	// Tunnel represents a port tunnel connection.
	Tunnel ConnectionType = "tunnel"
)

// ResourceConfig represents the configuration for connecting to an Azure resource.
type ResourceConfig struct {
	BastionHost    *config.BastionHost
	TargetResource *config.TargetResource
	Username       string
	LocalPort      int
	RemotePort     int
}
