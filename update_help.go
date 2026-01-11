package main

import (
	"strings"

	tea "github.com/charmbracelet/bubbletea"
)

// updateHelp handles keyboard input in the help view
func updateHelp(m Model, msg tea.KeyMsg) (Model, tea.Cmd) {
	// Calculate total lines and max scroll offset
	// We need to build the help content to know how many lines it has
	// This matches the logic in ui/help.go
	sections := []string{
		renderSymbolsSectionText(),
		"",
		renderCommandsSectionText(),
		"",
		renderNavigationSectionText(),
		"",
		renderDataStorageSectionText(),
	}
	content := strings.Join(sections, "\n")
	totalLines := strings.Count(content, "\n") + 1

	// Calculate available height (same as in RenderHelp)
	availableHeight := m.height - 3 // header + footer + message line
	if availableHeight < 5 {
		availableHeight = 5
	}

	maxOffset := totalLines - availableHeight
	if maxOffset < 0 {
		maxOffset = 0
	}

	// Check for update dismissal key 'u' (only if update notice is shown)
	if m.updateAvailable && !m.updateDismissed && msg.String() == "u" {
		m.updateDismissed = true
		return m, nil
	}

	switch msg.String() {
	case "q", "ctrl+c":
		return m, tea.Quit

	case "esc":
		// Return to previous view (where help was opened from)
		if m.previousView != "" {
			m.view = m.previousView
			m.previousView = ""
		} else {
			// Default to entries if no previous view
			m.view = "entries"
		}
		return m, nil

	case "f":
		// Scroll forward (down) one page
		m.scrollOffset += availableHeight
		if m.scrollOffset > maxOffset {
			m.scrollOffset = maxOffset
		}
		return m, nil

	case "b":
		// Scroll backward (up) one page
		m.scrollOffset -= availableHeight
		if m.scrollOffset < 0 {
			m.scrollOffset = 0
		}
		return m, nil

	default:
		return m, nil
	}
}

// Helper functions to build help content sections (text only, no styling)
func renderSymbolsSectionText() string {
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

func renderCommandsSectionText() string {
	title := "COMMANDS"

	items := []string{
		"Global:",
		"  n:new         Create new entry",
		"  a:todo        Add standalone todo",
		"  e:entries     Jump to entry list",
		"  t:todos       Jump to todo list",
		"  s:theme       Select theme",
		"  ?:help        Show this help page",
		"  q:quit        Quit application",
		"",
		"Entry List:",
		"  j/k:nav       Navigate up/down",
		"  enter:view    View entry detail",
		"  d:del         Mark/unmark for deletion",
		"  $:exec        Execute deletion (after marking)",
		"  /:filter      Filter by tags/dates",
		"  c:clear       Clear active filter",
		"",
		"Entry View:",
		"  j/k:nav       Navigate between entries",
		"  f/b:page      Page forward/backward in long entries",
		"  d:del         Mark/unmark for deletion",
		"  esc:back      Return to entry list",
		"",
		"Todo List:",
		"  j/k:nav       Navigate up/down",
		"  space:toggle  Cycle todo status (open → next → done)",
		"  d:del         Mark/unmark for deletion",
		"  $:exec        Execute deletion (after marking)",
		"  /:filter      Filter by tags/dates",
		"  c:clear       Clear active filter",
		"",
		"Forms (Entry/Todo):",
		"  ctrl+s:save   Save and close",
		"  esc:cancel    Cancel without saving",
		"",
		"Filter:",
		"  @tag          Filter by tag (tab to autocomplete)",
		"  today         Show items from today",
		"  yesterday     Show items from yesterday",
		"  last N days   Show items from last N days",
		"  YYYY-MM-DD    Show items from specific date",
		"  date to date  Show items in date range",
		"  enter:apply   Apply filter",
		"  esc:cancel    Cancel filter",
	}

	return title + "\n\n" + strings.Join(items, "\n")
}

func renderNavigationSectionText() string {
	title := "NAVIGATION"

	items := []string{
		"Entry list and todo list are peer views - use e/t to jump between them.",
		"Global shortcuts (n, a, s, ?) work from any read-only view.",
		"Press esc to exit forms and return to their natural home view.",
	}

	return title + "\n\n" + strings.Join(items, "\n")
}

func renderDataStorageSectionText() string {
	title := "DATA STORAGE"

	items := []string{
		"All data is stored in plain JSON files (default: ~/.amos/):",
		"  entries.json   Journal entries",
		"  todos.json     Todo items",
		"  config.json    User preferences (theme)",
		"",
		"Customize data directory for syncing with cloud storage:",
		"",
		"Option 1: Environment variable (one-time or in shell profile)",
		"  AMOS_DATA_DIR=~/Google\\ Drive/amos amos",
		"",
		"Option 2: Config file (persistent, create with any text editor)",
		"  Create file: ~/.config/amos/settings.json",
		"  Content:     {",
		"                 \"data_dir\": \"~/Google Drive/amos\"",
		"               }",
		"",
		"Priority: AMOS_DATA_DIR > ~/.config/amos/settings.json > ~/.amos",
		"Path formats: ~/path (tilde), /absolute/path, or ./relative/path",
	}

	return title + "\n\n" + strings.Join(items, "\n")
}
