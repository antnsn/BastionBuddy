package azure

import (
	"encoding/json"
	"fmt"
	"sync"

	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	"github.com/antnsn/BastionBuddy/internal/tunnels"
	"github.com/antnsn/BastionBuddy/internal/utils"
)

var (
	// Global state
	globalState struct {
		sync.RWMutex
		initialized    bool
		initializeOnce sync.Once
		cred           *azidentity.DefaultAzureCredential
		tunnelManager  *TunnelManager
	}
)

// checkAzLogin checks if the user is logged in to Azure CLI
func checkAzLogin() error {
	cmd := utils.PrepareAzureCommand("account", "show")
	output, err := cmd.CombinedOutput()
	if err != nil {
		// Try to parse the error to see if it's an authentication error
		var azError struct {
			Error struct {
				Code    string `json:"code"`
				Message string `json:"message"`
			} `json:"error"`
		}
		if jsonErr := json.Unmarshal(output, &azError); jsonErr == nil {
			if azError.Error.Code == "InvalidAuthenticationTokenTenant" ||
				azError.Error.Code == "ExpiredAuthenticationToken" {
				return fmt.Errorf("azure authentication required: %s", azError.Error.Message)
			}
		}
		return fmt.Errorf("failed to check Azure login status: %v", err)
	}
	return nil
}

// ensureAzLogin ensures the user is logged in to Azure CLI
func ensureAzLogin() error {
	if err := checkAzLogin(); err != nil {
		fmt.Println("Not logged in to Azure. Please follow the instructions to log in...")
		cmd := utils.PrepareAzureCommand("login")
		cmd.Stdout = nil // Let Azure CLI handle the output
		cmd.Stderr = nil
		cmd.Stdin = nil
		if err := cmd.Run(); err != nil {
			return fmt.Errorf("azure login failed: %v", err)
		}

		// Verify login was successful
		if err := checkAzLogin(); err != nil {
			return fmt.Errorf("login verification failed: %v", err)
		}
	}
	return nil
}

// initializeAzure sets up Azure credentials
func initializeAzure() error {
	// First ensure the user is logged in
	if err := ensureAzLogin(); err != nil {
		return err
	}

	// Then create the credential
	var err error
	globalState.cred, err = azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		return fmt.Errorf("failed to create Azure credential: %v", err)
	}
	return nil
}

// initializeTunnelManager initializes the tunnel manager
func initializeTunnelManager() error {
	configMgr, err := tunnels.NewManager()
	if err != nil {
		return fmt.Errorf("failed to initialize tunnel manager: %v", err)
	}

	globalState.tunnelManager = &TunnelManager{
		tunnels:   make(map[string]*TunnelInfo),
		configMgr: configMgr,
	}

	// Restore active tunnels from persistent storage
	activeTunnels := configMgr.GetActive()
	for _, t := range activeTunnels {
		tunnel := &TunnelInfo{
			ID:                    t.ID,
			LocalPort:             t.LocalPort,
			RemotePort:            t.RemotePort,
			ResourceID:            t.ResourceID,
			ResourceName:          t.ResourceName,
			SubscriptionID:        t.SubscriptionID,
			BastionName:           t.BastionName,
			BastionResourceGroup:  t.BastionResourceGroup,
			BastionSubscriptionID: t.BastionSubscriptionID,
			StartTime:             t.StartTime,
			Status:                t.Status,
			PID:                   t.PID,
		}
		globalState.tunnelManager.tunnels[t.ID] = tunnel
	}

	return nil
}

// initialize initializes all global state
func initialize() error {
	var initErr error
	globalState.initializeOnce.Do(func() {
		globalState.Lock()
		defer globalState.Unlock()

		// Initialize Azure credentials
		if err := initializeAzure(); err != nil {
			initErr = fmt.Errorf("azure initialization failed: %v", err)
			return
		}

		// Initialize tunnel manager
		if err := initializeTunnelManager(); err != nil {
			initErr = fmt.Errorf("tunnel manager initialization failed: %v", err)
			return
		}

		globalState.initialized = true
	})

	return initErr
}

// GetAzureCredential returns the Azure credential
func GetAzureCredential() (*azidentity.DefaultAzureCredential, error) {
	if err := initialize(); err != nil {
		return nil, err
	}

	globalState.RLock()
	defer globalState.RUnlock()
	return globalState.cred, nil
}

// GetTunnelManager returns the singleton instance of TunnelManager
func GetTunnelManager() (*TunnelManager, error) {
	if err := initialize(); err != nil {
		return nil, fmt.Errorf("failed to initialize: %v", err)
	}

	globalState.RLock()
	defer globalState.RUnlock()

	if !globalState.initialized {
		return nil, fmt.Errorf("global state not initialized")
	}

	return globalState.tunnelManager, nil
}

func init() {
	if err := initialize(); err != nil {
		fmt.Printf("Warning: initialization failed: %v\n", err)
	}
}
