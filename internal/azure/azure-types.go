// Package azure provides functionality for interacting with Azure resources
// through the Azure CLI, specifically for Bastion connections.
package azure

// ResourceConfig contains all the necessary configuration for connecting to an Azure resource.
type ResourceConfig struct {
	BastionHost    *BastionHost
	TargetResource *TargetResource
	Username       string
	LocalPort      int
	RemotePort     int
}
