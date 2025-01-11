package main

import (
	"fmt"
	"os"

	"github.com/yourusername/azbastion/internal/azure"
	"github.com/yourusername/azbastion/internal/config"
	"github.com/spf13/cobra"
)

func main() {
	cfg := config.NewConfig()
	
	rootCmd := &cobra.Command{
		Use:   "azbastion",
		Short: "Azure Bastion Connection Utility",
		Long: `A unified tool for managing Azure Bastion connections with support for SSH, RDP, and Tunneling.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := azure.CheckDependencies(); err != nil {
				return err
			}

			defer azure.Cleanup()

			if cfg.ConnectionType == "" {
				connType, err := azure.SelectConnectionType()
				if err != nil {
					return err
				}
				cfg.ConnectionType = connType
			}

			azureConfig, err := azure.GetAzureResources()
			if err != nil {
				return err
			}

			// Merge CLI config with Azure config
			azureConfig.Username = cfg.Username
			azureConfig.LocalPort = cfg.LocalPort
			azureConfig.RemotePort = cfg.RemotePort

			return azure.Connect(cfg.ConnectionType, azureConfig)
		},
	}

	// Add flags
	flags := rootCmd.Flags()
	flags.StringVarP(&cfg.ConnectionType, "type", "t", "", "Connection type (ssh|rdp|tunnel)")
	flags.StringVarP(&cfg.Username, "username", "u", "", "Username for SSH connection")
	flags.StringVarP(&cfg.LocalPort, "port", "p", "", "Local port for tunnel")
	flags.StringVarP(&cfg.RemotePort, "remote", "r", "", "Remote port for tunnel")

	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
