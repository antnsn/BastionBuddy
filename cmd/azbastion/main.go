// Package main is the entry point for BastionBuddy, a friendly command-line utility
// that makes Azure Bastion connections easy and interactive.
package main

import (
	"fmt"
	"os"

	"github.com/antnsn/BastionBuddy/internal/azure"
	"github.com/antnsn/BastionBuddy/internal/welcome"
)

// Version is the current version of BastionBuddy, set at build time.
var Version string

func main() {
	// Check for command line arguments
	if len(os.Args) > 1 {
		switch os.Args[1] {
		case "--tunnel":
			if len(os.Args) < 3 {
				fmt.Println("Error: Please specify a tunnel name")
				os.Exit(1)
			}
			tunnelName := os.Args[2]
			if err := azure.StartSavedTunnel(tunnelName); err != nil {
				fmt.Printf("Error starting tunnel: %v\n", err)
				os.Exit(1)
			}
			os.Exit(0)
		case "--ssh":
			if len(os.Args) < 3 {
				fmt.Println("Error: Please specify a saved configuration name")
				os.Exit(1)
			}
			configName := os.Args[2]
			if err := azure.StartSavedSSH(configName); err != nil {
				fmt.Printf("Error starting SSH: %v\n", err)
				os.Exit(1)
			}
			os.Exit(0)
		case "--rdp":
			if len(os.Args) < 3 {
				fmt.Println("Error: Please specify a saved configuration name")
				os.Exit(1)
			}
			configName := os.Args[2]
			if err := azure.StartSavedRDP(configName); err != nil {
				fmt.Printf("Error starting RDP: %v\n", err)
				os.Exit(1)
			}
			os.Exit(0)
		case "--list":
			var connectionType string
			if len(os.Args) > 2 {
				connectionType = os.Args[2]
			}
			if err := azure.ListConfigurations(connectionType); err != nil {
				fmt.Printf("Error listing configurations: %v\n", err)
				os.Exit(1)
			}
			os.Exit(0)
		}
	}

	welcome.ShowWelcome()

	if err := azure.CheckDependencies(); err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}

	defer func() {
		if err := azure.Cleanup(); err != nil {
			fmt.Printf("Error during cleanup: %v\n", err)
		}
	}()

	for {
		fmt.Println()
		action, err := azure.SelectInitialAction()
		if err != nil {
			fmt.Printf("Error: %v\n", err)
			continue
		}

		if action == "exit" {
			break
		}

		if err := azure.InitiateAction(action, nil); err != nil {
			fmt.Printf("Error: %v\n", err)
			continue
		}

		// Show welcome screen again with updated tunnel status
		fmt.Print("\033[H\033[2J") // Clear screen
		welcome.ShowWelcome()
	}
}
