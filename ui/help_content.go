package ui

import "strings"

// RenderSymbolsSectionText returns the plain text for the symbols section
func RenderSymbolsSectionText() string {
	title := "SYMBOLS"

	items := []string{
		"Entry List Markers:",
		"  +  Entry has linked todos (color shows highest priority)",
		"     Green = next, Magenta = open, Dim = done",
		"  D  Entry marked for deletion",
		"",
		"Todo List Markers:",
		"  +  Todo is linked to an entry",
		"  D  Todo marked for deletion",
		"",
		"Todo Status Indicators:",
		"  [ ] Open todo",
		"  [>] Next todo (high priority)",
		"  [x] Done todo",
	}

	return title + "\n\n" + strings.Join(items, "\n")
}

// RenderDataStorageSectionText returns the plain text for the data storage section
func RenderDataStorageSectionText() string {
	title := "DATA STORAGE"

	items := []string{
		"Data files stored in plain JSON (default: ~/.amos/):",
		"  entries.json   Journal entries",
		"  todos.json     Todo items",
		"",
		"Config stored in ~/.config/amos/settings.json:",
		"  data_dir       Custom data directory path",
		"  theme          Theme preference (brutalist, cyberpunk)",
		"",
		"Customize data directory for syncing with cloud storage:",
		"",
		"Option 1: Environment variable (one-time or in shell profile)",
		"  AMOS_DATA_DIR=~/Google\\ Drive/amos amos",
		"",
		"Option 2: Config file (persistent, create with any text editor)",
		"  {\"data_dir\": \"~/Google Drive/amos\", \"theme\": \"cyberpunk\"}",
		"",
		"Priority: AMOS_DATA_DIR > settings.json > ~/.amos",
		"Path formats: ~/path (tilde), /absolute/path, or ./relative/path",
	}

	return title + "\n\n" + strings.Join(items, "\n")
}
