package tunnels

import "time"

// Config represents a tunnel configuration
type Config struct {
	Name                  string    `json:"name"`
	SubscriptionID        string    `json:"subscription_id"`
	ResourceID            string    `json:"resource_id"`
	ResourceName          string    `json:"resource_name"`
	LocalPort             int       `json:"local_port"`
	RemotePort            int       `json:"remote_port"`
	Command               string    `json:"command"`
	Args                  []string  `json:"args"`
	LastUsed              time.Time `json:"last_used"`
	BastionName           string    `json:"bastion_name"`
	BastionResourceGroup  string    `json:"bastion_resource_group"`
	BastionSubscriptionID string    `json:"bastion_subscription_id"`
	ConnectionType        string    `json:"connection_type"`
	Username              string    `json:"username"`
	AuthType              string    `json:"auth_type"`
}

// SavedConfig represents a saved tunnel configuration
type SavedConfig struct {
	Name                  string    `json:"name"`
	LocalPort             int       `json:"local_port"`
	RemotePort            int       `json:"remote_port"`
	ResourceID            string    `json:"resource_id"`
	ResourceName          string    `json:"resource_name"`
	SubscriptionID        string    `json:"subscription_id"`
	LastUsed              time.Time `json:"last_used"`
	Command               string    `json:"command"`
	Args                  []string  `json:"args"`
	BastionName           string    `json:"bastion_name"`
	BastionResourceGroup  string    `json:"bastion_resource_group"`
	BastionSubscriptionID string    `json:"bastion_subscription_id"`
	ConnectionType        string    `json:"connection_type"`
	Username              string    `json:"username"`
}

// Active represents a currently running tunnel
type Active struct {
	ID                    string    `json:"id"`
	LocalPort             int       `json:"local_port"`
	RemotePort            int       `json:"remote_port"`
	ResourceID            string    `json:"resource_id"`
	ResourceName          string    `json:"resource_name"`
	SubscriptionID        string    `json:"subscription_id"`
	BastionName           string    `json:"bastion_name"`
	BastionResourceGroup  string    `json:"bastion_resource_group"`
	BastionSubscriptionID string    `json:"bastion_subscription_id"`
	StartTime             time.Time `json:"start_time"`
	Status                string    `json:"status"`
}
