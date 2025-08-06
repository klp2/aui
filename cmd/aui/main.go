package main

import (
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/yourusername/aui/internal/ui"
)

func main() {
	// Create the initial app state
	app := ui.InitialApp()

	// Create the Bubble Tea program
	p := tea.NewProgram(app, tea.WithAltScreen())

	// Run the program
	if _, err := p.Run(); err != nil {
		fmt.Printf("Error running aui: %v\n", err)
		os.Exit(1)
	}
}
