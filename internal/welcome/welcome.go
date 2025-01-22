// Package welcome provides functionality for displaying the welcome screen
// and initial application information to users.
package welcome

import (
	"fmt"
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
	logo := `
 ██████╗  █████╗ ███████╗████████╗██╗ ██████╗ ███╗   ██╗██████╗ ██╗   ██╗██████╗ ██████╗ ██╗   ██╗
 ██╔══██╗██╔══██╗██╔════╝╚══██╔══╝██║██╔═══██╗████╗  ██║██╔══██╗██║   ██║██╔══██╗██╔══██╗╚██╗ ██╔╝
 ██████╔╝███████║███████╗   ██║   ██║██║   ██║██╔██╗ ██║██████╔╝██║   ██║██║  ██║██║  ██║ ╚████╔╝
 ██╔══██╗██╔══██║╚════██║   ██║   ██║██║   ██║██║╚██╗██║██╔══██╗██║   ██║██║  ██║██║  ██║  ╚██╔╝
 ██████╔╝██║  ██║███████║   ██║   ██║╚██████╔╝██║ ╚████║██████╔╝╚██████╔╝██████╔╝██████╔╝   ██║
 ╚═════╝ ╚═╝  ╚═╝╚══════╝   ╚═╝   ╚═╝ ╚═════╝ ╚═╝  ╚═══╝╚═════╝  ╚═════╝ ╚═════╝ ╚═════╝    ╚═╝
`
	if _, err := cyan.Println(logo); err != nil {
		fmt.Println(logo) // Fallback to regular print if colored fails
	}

	if _, err := magenta.Print("Version: "); err != nil {
		fmt.Print("Version: ")
	}
	fmt.Println(Version)

	if _, err := magenta.Print("Description: "); err != nil {
		fmt.Print("Description: ")
	}
	fmt.Println("A friendly command-line utility that makes Azure Bastion connections easy and interactive.")

	fmt.Println()
	if _, err := yellow.Println("✨ Features:"); err != nil {
		fmt.Println("✨ Features:")
	}
	fmt.Println("• Interactive menu-driven interface")
	fmt.Println("• Support for SSH, RDP, and Port Tunneling")
	fmt.Println("• Automatic Azure resource discovery")
	fmt.Println("• Kill active tunnels with ease")

	fmt.Println()
	if _, err := yellow.Println("🚀 Usage Tips:"); err != nil {
		fmt.Println("🚀 Usage Tips:")
	}
	fmt.Println("• Use arrow keys to navigate")
	fmt.Println("• Type to search in lists")
	fmt.Println("• Press Enter to select")
	fmt.Println("• Use Ctrl+C to exit at any time")
	fmt.Println("• Select 'Manage Tunnels' to manage tunnel connections")

	printSeparator()
	showActiveTunnels() // Show active tunnels at the bottom of the welcome screen
	printSeparator()
}

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
		if t.Status == "restored" {
			status = "Restored"
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
