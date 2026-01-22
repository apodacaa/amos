package main

import (
	"strings"

	"github.com/apodacaa/amos/ui"
	tea "github.com/charmbracelet/bubbletea"
)

// updateHelp handles keyboard input in the help view
func updateHelp(m Model, msg tea.KeyMsg) (Model, tea.Cmd) {
	// Calculate total lines and max scroll offset
	// We need to build the help content to know how many lines it has
	// This matches the logic in ui/help.go
	content := ui.GetHelpContent()
	totalLines := len(strings.Split(content, "\n"))

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
