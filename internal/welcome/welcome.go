package welcome

import (
	"fmt"
	"strings"

	"github.com/fatih/color"
)

const logo = `    ____             __  _             ____            __    __     
   / __ )____ ______/ /_(_)___  ____  / __ )__  ______/ /___/ /_  __
  / __  / __ '/ ___/ __/ / __ \/ __ \/ __  / / / / __  / __  / / / /
 / /_/ / /_/ (__  ) /_/ / /_/ / / / / /_/ / /_/ / /_/ / /_/ / /_/ / 
/_____/\__,_/____/\__/_/\____/_/ /_/_____/\__,_/\__,_/\__,_/\__, /  
                                                           /____/   `

func ShowWelcome() {
	// Clear the screen first
	fmt.Print("\033[H\033[2J")

	// Print the logo in bright cyan
	cyan := color.New(color.FgHiCyan, color.Bold)
	cyan.Println(logo)

	// Print version and description
	magenta := color.New(color.FgHiMagenta)
	magenta.Print("Version: ")
	fmt.Println("1.0.0")
	magenta.Print("Description: ")
	fmt.Println("A friendly Azure Bastion connection utility")
	fmt.Println()

	// Print features
	yellow := color.New(color.FgYellow, color.Bold)
	yellow.Println("✨ Features:")
	features := []string{
		"🔒 Secure SSH connections to Azure VMs",
		"🌐 Port tunneling for remote access",
		"⚡ Smart caching for faster resource listing",
		"🎯 Interactive menu navigation",
	}
	for _, feature := range features {
		fmt.Printf("   %s\n", feature)
	}
	fmt.Println()

	// Print usage tips
	yellow.Println("🚀 Usage Tips:")
	tips := []string{
		"↑/↓  Navigate through options",
		"⌨️   Type to search in lists",
		"↵    Press Enter to select",
		"^C   Press Ctrl+C to exit",
	}
	for _, tip := range tips {
		fmt.Printf("   %s\n", tip)
	}
	fmt.Println()

	// Print separator with gradient effect
	colors := []*color.Color{
		color.New(color.FgHiBlue),
		color.New(color.FgHiCyan),
		color.New(color.FgHiMagenta),
	}
	
	separator := strings.Repeat("─", 60)
	parts := len(colors)
	partLength := len(separator) / parts
	
	for i, c := range colors {
		start := i * partLength
		end := start + partLength
		if i == len(colors)-1 {
			end = len(separator)
		}
		c.Print(separator[start:end])
	}
	fmt.Println("\n")
}
