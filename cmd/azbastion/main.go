// Package main provides the entry point for the BastionBuddy CLI tool,
// which facilitates SSH connections and tunneling to Azure virtual machines
// through Azure Bastion.
package main

import (
	"fmt"
	"os"

	"github.com/antnsn/BastionBuddy/internal/azure"
	"github.com/antnsn/BastionBuddy/internal/welcome"
	"github.com/spf13/cobra"
)

func main() {
	rootCmd := &cobra.Command{
		Use:   "bastionbuddy",
		Short: "Azure Bastion Connection Utility",
		Long:  "A utility for managing connections to Azure VMs through Azure Bastion",
		RunE: func(_ *cobra.Command, _ []string) error {
			welcome.ShowWelcome()

			if err := azure.CheckDependencies(); err != nil {
				return fmt.Errorf("dependency check failed: %v", err)
			}
			defer func() {
				if err := azure.Cleanup(); err != nil {
					fmt.Fprintf(os.Stderr, "Error during cleanup: %v\n", err)
				}
			}()

			config, err := azure.GetAzureResources()
			if err != nil {
				return fmt.Errorf("failed to get Azure resources: %v", err)
			}

			connectionType, err := azure.SelectConnectionType()
			if err != nil {
				return fmt.Errorf("failed to select connection type: %v", err)
			}

			if err := azure.Connect(connectionType, config); err != nil {
				return fmt.Errorf("failed to establish connection: %v", err)
			}

			return nil
		},
	}

	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}
