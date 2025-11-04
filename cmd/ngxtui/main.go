package main

import (
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"

	"github.com/aitmiloud/ngxtui/internal/app"
)

func main() {
	// Check if running as root
	if os.Geteuid() != 0 {
		fmt.Println("This program must be run as root (use sudo)")
		os.Exit(1)
	}

	// Initialize the model
	initialModel := app.InitialModel()

	// Create the program with alt screen and mouse support
	p := tea.NewProgram(
		&app.App{Model: initialModel},
		tea.WithAltScreen(),
		tea.WithMouseCellMotion(),
	)

	// Run the program
	if _, err := p.Run(); err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}
}
