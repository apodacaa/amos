package main

import (
	"fmt"
	"os"

	"github.com/apodacaa/amos/cli"
	"github.com/apodacaa/amos/internal/helpers"
	"github.com/apodacaa/amos/internal/storage"
	tea "github.com/charmbracelet/bubbletea"
)

// Version is set via ldflags during build
var Version = "dev"

func main() {
	// Handle --help flag
	if len(os.Args) > 1 && (os.Args[1] == "--help" || os.Args[1] == "-h" || os.Args[1] == "help") {
		fmt.Printf("amos %s - journal + todo manager\n\n", Version)
		fmt.Println("Usage:")
		fmt.Println("  amos                                    Launch TUI")
		fmt.Println("  amos add entry --title T [--body B]     Add a journal entry")
		fmt.Println("  amos add todo --title T [--status S]    Add a todo (status: open|next|done)")
		fmt.Println("  amos list entries [--tag TAG]            List entries as JSON")
		fmt.Println("  amos list todos [--tag TAG] [--status S] List todos as JSON")
		fmt.Println("  amos update todo <id> --status S        Update a todo's status (open|next|done)")
		fmt.Println("  amos done <id>                          Mark a todo done (alias for update todo --status done)")
		fmt.Println("")
		fmt.Println("Flags:")
		fmt.Println("  -v, --version       Print version")
		fmt.Println("  --check-update      Check for updates")
		fmt.Println("  -h, --help          Show this help")
		os.Exit(0)
	}

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
	// Must happen before NewModel() or CLI commands which load data
	if err := storage.InitializeDataDir(); err != nil {
		fmt.Fprintf(os.Stderr, "Error initializing data directory: %v\n", err)
		os.Exit(1)
	}

	// Handle CLI subcommands (add, list)
	if cli.Run(os.Args) {
		return
	}

	m := NewModel()
	p := tea.NewProgram(m, tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}
