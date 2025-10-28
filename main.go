package main

import (
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"
)

// Version is set via ldflags during build
var Version = "dev"

func main() {
	// Handle --version flag
	if len(os.Args) > 1 && (os.Args[1] == "--version" || os.Args[1] == "-v") {
		fmt.Printf("amos version %s\n", Version)
		os.Exit(0)
	}

	m := NewModel()
	p := tea.NewProgram(m, tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}
