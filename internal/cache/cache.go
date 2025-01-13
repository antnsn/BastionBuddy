// Package cache provides caching functionality for Azure resource queries
// to improve performance of subsequent requests.
package cache

import (
	"fmt"
	"os"
	"path/filepath"
	"time"
)

// Cache represents a file-based caching system for storing Azure resource query results.
type Cache struct {
	dir     string        // Directory where cache files are stored
	timeout time.Duration // How long cache entries are valid
}

// NewCache creates a new Cache instance with the specified directory and timeout duration.
func NewCache(dir string, timeout time.Duration) *Cache {
	return &Cache{
		dir:     dir,
		timeout: timeout,
	}
}

// Get retrieves a value from the cache. If the value doesn't exist or is expired,
// it calls the fetch function to get and cache a new value.
func (c *Cache) Get(key string, fetch func() ([]byte, error)) ([]byte, error) {
	path := filepath.Join(c.dir, key)

	// Check if cache file exists and is not expired
	if stat, err := os.Stat(path); err == nil {
		if time.Since(stat.ModTime()) < c.timeout {
			// Cache hit
			return os.ReadFile(path)
		}
		// File exists but is expired, remove it
		if err := os.Remove(path); err != nil {
			return nil, fmt.Errorf("failed to remove expired cache file: %v", err)
		}
	}

	// Cache miss or expired, fetch new data
	data, err := fetch()
	if err != nil {
		return nil, err
	}

	// Ensure cache directory exists
	if err := os.MkdirAll(c.dir, 0755); err != nil {
		return nil, err
	}

	// Write to cache
	if err := os.WriteFile(path, data, 0644); err != nil {
		return nil, fmt.Errorf("failed to write to cache: %v", err)
	}

	return data, nil
}

// Cleanup removes all expired cache files from the cache directory.
func (c *Cache) Cleanup() error {
	return filepath.Walk(c.dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return nil
		}
		if !info.IsDir() && time.Since(info.ModTime()) > 24*time.Hour {
			if err := os.Remove(path); err != nil {
				return fmt.Errorf("failed to remove expired cache file: %v", err)
			}
		}
		return nil
	})
}
