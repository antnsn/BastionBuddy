// Package welcome provides functionality for displaying the welcome screen
// and initial application information to users.
package welcome

import (
	"fmt"
	"strings"

	"github.com/fatih/color"
)

var (
	cyan    = color.New(color.FgCyan)
	magenta = color.New(color.FgMagenta)
	yellow  = color.New(color.FgYellow)
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
	fmt.Println("1.0.2")

	if _, err := magenta.Print("Description: "); err != nil {
		fmt.Print("Description: ")
	}
	fmt.Println("A friendly command-line utility that makes Azure Bastion connections easy and interactive.")

	fmt.Println()
	if _, err := yellow.Println("✨ Features:"); err != nil {
		fmt.Println("✨ Features:")
	}
	fmt.Println("• Interactive menu-driven interface")
	fmt.Println("• Support for both SSH and Port Tunneling")
	fmt.Println("• Automatic Azure resource discovery")
	fmt.Println("• Smart caching for faster subsequent connections")
	fmt.Println("• Colorful and intuitive UI")

	fmt.Println()
	if _, err := yellow.Println("🚀 Usage Tips:"); err != nil {
		fmt.Println("🚀 Usage Tips:")
	}
	fmt.Println("• Use arrow keys to navigate")
	fmt.Println("• Type to search in lists")
	fmt.Println("• Press Enter to select")
	fmt.Println("• Use Ctrl+C to exit at any time")

	printSeparator()
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
