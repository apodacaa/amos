package main

import (
	"fmt"
	"os"

	"github.com/apodacaa/amos/internal/helpers"
	"github.com/apodacaa/amos/internal/storage"
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

	// Handle --check-update flag
	if len(os.Args) > 1 && os.Args[1] == "--check-update" {
		fmt.Printf("Current version: %s\n", Version)
		fmt.Println("Checking for updates...")

		latest := helpers.FetchLatestVersion()
		if latest == "" {
			fmt.Println("✗ Failed to fetch latest version from GitHub")
			fmt.Println("  Possible causes:")
			fmt.Println("  - Network connectivity issue")
			fmt.Println("  - GitHub API rate limit (60 requests/hour)")
			fmt.Println("  - API timeout (3 second limit)")
			fmt.Println("  - GitHub service issue")
			os.Exit(1)
		}

		fmt.Printf("Latest version: %s\n", latest)

		if helpers.IsUpdateAvailable(Version, latest) {
			fmt.Printf("✓ Update available: %s → %s\n", Version, latest)
			fmt.Println("\nInstall with:")
			fmt.Println("  brew upgrade amos")
			fmt.Println("  choco upgrade amos")
			fmt.Println("  Or visit: https://github.com/apodacaa/amos/releases")
		} else {
			fmt.Println("✓ You are running the latest version")
		}
		return
	}

	// Initialize data directory from env var or system config
	// Must happen before NewModel() which loads data
	if err := storage.InitializeDataDir(); err != nil {
		fmt.Fprintf(os.Stderr, "Error initializing data directory: %v\n", err)
		os.Exit(1)
	}

	m := NewModel()
	p := tea.NewProgram(m, tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}
