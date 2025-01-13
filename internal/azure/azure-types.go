package azure

// Resource represents an Azure resource
type Resource struct {
	Name string `json:"name"`
	ID   string `json:"id"`
}

// BastionHost represents an Azure Bastion host
type BastionHost struct {
	Name          string `json:"name"`
	ResourceGroup string `json:"resourceGroup"`
}

// AzureConfig holds the configuration for Azure resources
type AzureConfig struct {
	ResourceID            string
	BastionName          string
	BastionResourceGroup string
	Username             string
	LocalPort            int
	RemotePort           int
}
