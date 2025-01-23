// Package config provides configuration functionality for BastionBuddy.
package config

// Config represents the application configuration.
type Config struct {
	Username       string
	ResourceGroup  string
	SubscriptionID string
	ConnectionType string
	LocalPort      int
	RemotePort     int
}

// NewConfig creates a new Config instance with default values.
func NewConfig() *Config {
	return &Config{
		LocalPort:  22,
		RemotePort: 22,
	}
}
