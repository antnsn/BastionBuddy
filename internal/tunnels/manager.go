// Package tunnels provides functionality for managing tunnel connections
// and their configurations.
package tunnels

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

// Manager handles saving and loading tunnel configurations
type Manager struct {
	configDir     string
	tunnelFile    string
	sshFile       string
	rdpFile       string
	activeTunnels string
	tunnelConfigs []Config
	sshConfigs    []Config
	rdpConfigs    []Config
	active        []Active
}

// NewManager creates a new tunnel configuration manager
func NewManager() (*Manager, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return nil, fmt.Errorf("failed to get user home directory: %v", err)
	}

	configDir := filepath.Join(homeDir, ".config", "bastionbuddy")
	if err := os.MkdirAll(configDir, 0755); err != nil {
		return nil, fmt.Errorf("failed to create config directory: %v", err)
	}

	manager := &Manager{
		configDir:     configDir,
		tunnelFile:    filepath.Join(configDir, "tunnels.json"),
		sshFile:       filepath.Join(configDir, "ssh.json"),
		rdpFile:       filepath.Join(configDir, "rdp.json"),
		activeTunnels: filepath.Join(configDir, "active.json"),
	}

	if err := manager.load(); err != nil {
		return nil, err
	}

	return manager, nil
}

// SaveConfig saves a tunnel configuration for future use
func (m *Manager) SaveConfig(config Config) error {
	switch config.ConnectionType {
	case "ssh":
		// Check if configuration with same name exists
		for i, existing := range m.sshConfigs {
			if existing.Name == config.Name {
				// Update existing configuration
				m.sshConfigs[i] = config
				return m.save()
			}
		}
		// Add new configuration if not found
		m.sshConfigs = append(m.sshConfigs, config)
	case "rdp":
		// Check if configuration with same name exists
		for i, existing := range m.rdpConfigs {
			if existing.Name == config.Name {
				// Update existing configuration
				m.rdpConfigs[i] = config
				return m.save()
			}
		}
		// Add new configuration if not found
		m.rdpConfigs = append(m.rdpConfigs, config)
	default:
		// Check if configuration with same name exists
		for i, existing := range m.tunnelConfigs {
			if existing.Name == config.Name {
				// Update existing configuration
				m.tunnelConfigs[i] = config
				return m.save()
			}
		}
		// Add new configuration if not found
		m.tunnelConfigs = append(m.tunnelConfigs, config)
	}
	return m.save()
}

// GetSavedConfigs returns all saved tunnel configurations
func (m *Manager) GetSavedConfigs() []Config {
	return append(append(m.tunnelConfigs, m.sshConfigs...), m.rdpConfigs...)
}

// GetSavedConfigsByType returns saved configurations of a specific type
func (m *Manager) GetSavedConfigsByType(connectionType string) []Config {
	switch connectionType {
	case "ssh":
		return m.sshConfigs
	case "rdp":
		return m.rdpConfigs
	default:
		return m.tunnelConfigs
	}
}

// SaveActive saves information about a currently active tunnel
func (m *Manager) SaveActive(tunnel Active) error {
	m.active = append(m.active, tunnel)
	return m.save()
}

// RemoveActive removes a tunnel from the active tunnels list
func (m *Manager) RemoveActive(id string) error {
	for i, t := range m.active {
		if t.ID == id {
			m.active = append(m.active[:i], m.active[i+1:]...)
			return m.save()
		}
	}
	return nil
}

// GetActive returns all currently active tunnels
func (m *Manager) GetActive() []Active {
	return m.active
}

// load loads the saved configurations and active tunnels from disk
func (m *Manager) load() error {
	// Load tunnel configurations
	if data, err := os.ReadFile(m.tunnelFile); err == nil {
		if err := json.Unmarshal(data, &m.tunnelConfigs); err != nil {
			return fmt.Errorf("failed to parse tunnel configurations: %v", err)
		}
	}

	// Load SSH configurations
	if data, err := os.ReadFile(m.sshFile); err == nil {
		if err := json.Unmarshal(data, &m.sshConfigs); err != nil {
			return fmt.Errorf("failed to parse SSH configurations: %v", err)
		}
	}

	// Load RDP configurations
	if data, err := os.ReadFile(m.rdpFile); err == nil {
		if err := json.Unmarshal(data, &m.rdpConfigs); err != nil {
			return fmt.Errorf("failed to parse RDP configurations: %v", err)
		}
	}

	// Load active tunnels
	if data, err := os.ReadFile(m.activeTunnels); err == nil {
		if err := json.Unmarshal(data, &m.active); err != nil {
			return fmt.Errorf("failed to parse active tunnels: %v", err)
		}
	}

	return nil
}

// save saves the current configurations and active tunnels to disk
func (m *Manager) save() error {
	// Save tunnel configurations
	data, err := json.MarshalIndent(m.tunnelConfigs, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal tunnel configurations: %v", err)
	}
	if err := os.WriteFile(m.tunnelFile, data, 0644); err != nil {
		return fmt.Errorf("failed to save tunnel configurations: %v", err)
	}

	// Save SSH configurations
	data, err = json.MarshalIndent(m.sshConfigs, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal SSH configurations: %v", err)
	}
	if err := os.WriteFile(m.sshFile, data, 0644); err != nil {
		return fmt.Errorf("failed to save SSH configurations: %v", err)
	}

	// Save RDP configurations
	data, err = json.MarshalIndent(m.rdpConfigs, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal RDP configurations: %v", err)
	}
	if err := os.WriteFile(m.rdpFile, data, 0644); err != nil {
		return fmt.Errorf("failed to save RDP configurations: %v", err)
	}

	// Save active tunnels
	data, err = json.MarshalIndent(m.active, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal active tunnels: %v", err)
	}
	if err := os.WriteFile(m.activeTunnels, data, 0644); err != nil {
		return fmt.Errorf("failed to save active tunnels: %v", err)
	}

	return nil
}
