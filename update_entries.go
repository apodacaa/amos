package main

import (
	"time"

	"github.com/apodacaa/amos/internal/helpers"
	tea "github.com/charmbracelet/bubbletea"
)

// handleEntriesListKeys processes keyboard input (entries list view)
func (m Model) handleEntriesListKeys(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "q", "ctrl+c":
		return m, tea.Quit
	case "esc":
		// Go back to dashboard
		m.view = "dashboard"
		m.statusMsg = "" // Clear status message when changing views
		return m, nil
	case "n":
		// Create new entry (using shared helper)
		return m.handleNewEntry()
	case "t":
		// Jump to todo list (explicit navigation)
		m.view = "todos"
		m.selectedTodo = 0
		m.statusMsg = "" // Clear status message when changing views
		// Entries already loaded, just need to ensure todos are loaded
		// (but loadEntriesAndTodos is safe to call again)
		return m, m.loadEntriesAndTodos()
	case "a":
		// Add standalone todo (using shared helper)
		return m.handleAddTodo()
	case "@":
		// Open tag picker (or clear filter if already filtering)
		if m.filterTag != "" {
			// Clear filter
			m.filterTag = ""
			m.statusMsg = "✓ Filter cleared"
			m.statusTime = time.Now()
			return m, clearStatusAfterDelay()
		}
		// Extract unique tags from all entries
		m.availableTags = helpers.ExtractUniqueTags(m.entries)
		if len(m.availableTags) == 0 {
			m.statusMsg = "No tags found in entries"
			m.statusTime = time.Now()
			return m, clearStatusAfterDelay()
		}
		m.selectedTag = 0
		m.view = "tag_picker"
		return m, nil
	case "j", "down":
		if m.selectedEntry < len(m.entries)-1 {
			m.selectedEntry++
		}
		return m, nil
	case "k", "up":
		if m.selectedEntry > 0 {
			m.selectedEntry--
		}
		return m, nil
	case "enter":
		// Open selected entry for read-only viewing
		// Apply filter if active (same logic as UI)
		filtered := helpers.FilterEntriesByTag(m.entries, m.filterTag)

		if m.selectedEntry >= 0 && m.selectedEntry < len(filtered) {
			// Need to get the sorted entry (newest first)
			sorted := helpers.SortEntriesForDisplay(filtered)
			m.viewingEntry = sorted[m.selectedEntry]
			m.view = "view_entry"
			// Load todos so we can display them in the entry view
			return m, m.loadTodos()
		}
		return m, nil
	}
	return m, nil
}
