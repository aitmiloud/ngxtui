package main

import (
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"

	"github.com/aitmiloud/ngxtui/internal/app"
)

func main() {
	// Check if running with appropriate permissions
	if os.Geteuid() != 0 {
		fmt.Println("⚠️  Warning: This application requires root privileges to manage NGINX.")
		fmt.Println("Please run with sudo: sudo ngxtui")
		os.Exit(1)
	}

	// Create app instance
	a := app.New()

	// Create Bubble Tea program
	p := tea.NewProgram(a, tea.WithAltScreen())

	// Run the program
	if _, err := p.Run(); err != nil {
		fmt.Printf("Error running program: %v\n", err)
		os.Exit(1)
	}
}
