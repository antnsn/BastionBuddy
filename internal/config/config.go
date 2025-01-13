// Package config provides configuration functionality for BastionBuddy.
package config

import (
	"os"
	"path/filepath"
	"time"
)

// Config represents the application configuration.
type Config struct {
	CacheDir       string
	Username       string
	ResourceGroup  string
	SubscriptionID string
	ConnectionType string
	LocalPort      int
	RemotePort     int
	CacheTimeout   time.Duration
}

// NewConfig creates a new Config instance with default values.
func NewConfig() *Config {
	return &Config{
		CacheDir:     defaultCacheDir(),
		CacheTimeout: 24 * time.Hour,
	}
}

// defaultCacheDir returns the default cache directory path.
func defaultCacheDir() string {
	cacheDir := os.TempDir()
	if cacheDir == "" {
		cacheDir = os.Getenv("HOME")
		if cacheDir == "" {
			cacheDir = "."
		}
	}
	return filepath.Join(cacheDir, "bastionbuddy-cache")
}
