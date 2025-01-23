package config

// BastionHost represents an Azure Bastion host
type BastionHost struct {
	Name           string
	ResourceGroup  string
	SubscriptionID string
}

// TargetResource represents an Azure resource that can be connected to via Bastion
type TargetResource struct {
	ID             string
	Name           string
	Type           string
	SubscriptionID string
}

// ResourceConfig represents the configuration for connecting to an Azure resource
type ResourceConfig struct {
	BastionHost    *BastionHost
	TargetResource *TargetResource
	Username       string
	LocalPort      int
	RemotePort     int
}
