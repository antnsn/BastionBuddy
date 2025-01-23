// Package welcome provides functionality for displaying the welcome screen
// and initial application information to users.
package welcome

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"time"

	"github.com/antnsn/BastionBuddy/internal/azure"
	"github.com/fatih/color"
)

var (
	// Version holds the current version of BastionBuddy
	// This will be set during build time using ldflags
	Version = "dev"

	cyan    = color.New(color.FgCyan)
	magenta = color.New(color.FgMagenta)
	yellow  = color.New(color.FgYellow)
	green   = color.New(color.FgGreen)
)

// ShowWelcome displays the welcome screen with the application logo,
// version information, and usage tips.
func ShowWelcome() {
	// Clear screen first
	fmt.Print("\033[H\033[2J")

	logo := `
 ██████╗  █████╗ ███████╗████████╗██╗ ██████╗ ███╗   ██╗██████╗ ██╗   ██╗██████╗ ██████╗ ██╗   ██╗
 ██╔══██╗██╔══██╗██╔════╝╚══██╔══╝██║██╔═══██╗████╗  ██║██╔══██╗██║   ██║██╔══██╗██╔══██╗╚██╗ ██╔╝
 ██████╔╝███████║███████╗   ██║   ██║██║   ██║██╔██╗ ██║██████╔╝██║   ██║██║  ██║██║  ██║ ╚████╔╝
 ██╔══██╗██╔══██║╚════██║   ██║   ██║██║   ██║██║╚██╗██║██╔══██╗██║   ██║██║  ██║██║  ██║  ╚██╔╝
 ██████╔╝██║  ██║███████║   ██║   ██║╚██████╔╝██║ ╚████║██████╔╝╚██████╔╝██████╔╝██████╔╝   ██║
 ╚═════╝ ╚═╝  ╚═╝╚══════╝   ╚═╝   ╚═╝ ╚═════╝ ╚═╝  ╚═══╝╚═════╝  ╚═════╝ ╚═════╝ ╚═════╝    ╚═╝
`
	// Print logo
	if _, err := cyan.Println(logo); err != nil {
		fmt.Println(logo)
	}

	// Print version
	if _, err := magenta.Print("Version: "); err != nil {
		fmt.Print("Version: ")
	}
	fmt.Println(Version)

	if _, err := magenta.Print("🌟 Description: "); err != nil {
		fmt.Print("🌟 Description: ")
	}
	fmt.Println("Your friendly Azure Bastion companion for seamless cloud connections")

	fmt.Println()
	if _, err := yellow.Println("✨ Features & Capabilities:"); err != nil {
		fmt.Println("✨ Features & Capabilities:")
	}
	fmt.Println("  🔒 Secure SSH connections with saved configurations")
	fmt.Println("  🖥️ Remote Desktop (RDP) support for Windows VMs")
	fmt.Println("  🌐 Port tunneling with connection management")
	fmt.Println("  🎯 Smart resource discovery and caching")
	fmt.Println("  💾 Save and reuse your favorite connections")

	fmt.Println()
	if _, err := yellow.Println("⚡ Quick Commands:"); err != nil {
		fmt.Println("⚡ Quick Commands:")
	}
	fmt.Println("  • bastionbuddy list           → View saved configs")
	fmt.Println("  • bastionbuddy ssh <name>     → Quick SSH connection")
	fmt.Println("  • bastionbuddy rdp <name>     → Start RDP session")
	fmt.Println("  • bastionbuddy tunnel <name>  → Create port tunnel")

	fmt.Println()
	if _, err := yellow.Println("🎮 Navigation Tips:"); err != nil {
		fmt.Println("🎮 Navigation Tips:")
	}
	fmt.Println("  ↑↓ Arrow keys to navigate menus")
	fmt.Println("  ⌨️  Type to filter and search")
	fmt.Println("  ⏎  Enter to select")

	fmt.Println()
	if _, err := cyan.Print("📂 Config Location: "); err != nil {
		fmt.Print("📂 Config Location: ")
	}
	if homeDir, err := os.UserHomeDir(); err == nil {
		configPath := filepath.Join(homeDir, ".config", "bastionbuddy")
		if runtime.GOOS == "windows" {
			// Convert to Windows path style for display
			configPath = strings.ReplaceAll(configPath, "/", "\\")
		}
		fmt.Printf("%s\n", configPath)
	} else {
		fmt.Printf("~/.config/bastionbuddy/\n")
	}

	printSeparator()

	// Add an empty line before active tunnels for dynamic updates
	fmt.Println()
	showActiveTunnels()

	printSeparator()
}

// showActiveTunnels displays the list of active tunnels
func showActiveTunnels() {
	manager, err := azure.GetTunnelManager()
	if err != nil {
		fmt.Printf("Warning: failed to get tunnel manager: %v\n", err)
		return
	}

	tunnels := manager.ListTunnels()
	if len(tunnels) == 0 {
		return
	}

	if _, err := yellow.Println("🔌 Active Tunnels:"); err != nil {
		fmt.Println("🔌 Active Tunnels:")
	}

	for _, t := range tunnels {
		status := "Running"
		if t.Status != "running" {
			status = "Unknown"
		}

		if _, err := green.Printf("• %s (Local:%d → Remote:%d) - Resource: %s [%s: %s]\n",
			t.ID[:8], t.LocalPort, t.RemotePort, t.ResourceName,
			status, time.Since(t.StartTime).Round(time.Second)); err != nil {
			fmt.Printf("• %s (Local:%d → Remote:%d) - Resource: %s [%s: %s]\n",
				t.ID[:8], t.LocalPort, t.RemotePort, t.ResourceName,
				status, time.Since(t.StartTime).Round(time.Second))
		}
	}
}

func printSeparator() {
	separator := strings.Repeat("=", 80)
	colors := []*color.Color{
		color.New(color.FgBlue),
		color.New(color.FgMagenta),
		color.New(color.FgCyan),
	}

	partLength := len(separator) / len(colors)

	for i, c := range colors {
		start := i * partLength
		end := start + partLength
		if i == len(colors)-1 {
			end = len(separator)
		}
		if _, err := c.Print(separator[start:end]); err != nil {
			fmt.Print(separator[start:end]) // Fallback to regular print if colored fails
		}
	}
	fmt.Print("\n\n")
}
