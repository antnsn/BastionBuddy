package cache

import (
	"os"
	"path/filepath"
	"time"
)

type Cache struct {
	dir     string
	timeout time.Duration
}

func NewCache(dir string, timeout time.Duration) *Cache {
	return &Cache{
		dir:     dir,
		timeout: timeout,
	}
}

func (c *Cache) Get(key string, fetch func() ([]byte, error)) ([]byte, error) {
	path := filepath.Join(c.dir, key)
	
	if err := os.MkdirAll(c.dir, 0755); err != nil {
		return nil, err
	}

	if info, err := os.Stat(path); err == nil {
		if time.Since(info.ModTime()) < c.timeout {
			return os.ReadFile(path)
		}
	}

	data, err := fetch()
	if err != nil {
		return nil, err
	}

	if err := os.WriteFile(path, data, 0644); err != nil {
		return nil, err
	}

	return data, nil
}

func (c *Cache) Cleanup() error {
	return filepath.Walk(c.dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return nil
		}
		if !info.IsDir() && time.Since(info.ModTime()) > 24*time.Hour {
			os.Remove(path)
		}
		return nil
	})
}
