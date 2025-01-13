package config

import (
	"os"
	"path/filepath"
)

type Config struct {
	ConnectionType string
	Username      string
	LocalPort     int
	RemotePort    int
	CacheTimeout  int
	CacheDir      string
}

func NewConfig() *Config {
	return &Config{
		CacheTimeout: 3600, // 1 hour in seconds
		CacheDir:    defaultCacheDir(),
	}
}

func defaultCacheDir() string {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return filepath.Join(os.TempDir(), ".cache", "azbastion")
	}
	return filepath.Join(homeDir, ".cache", "azbastion")
}
