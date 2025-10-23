package main

import (
	"time"

	"github.com/apodacaa/amos/internal/helpers"
	"github.com/charmbracelet/bubbles/textarea"
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
		// Open tag filter input (or clear filter if already filtering)
		if len(m.filterTags) > 0 {
			// Clear filters
			m.filterTags = []string{}
			m.statusMsg = "✓ Filter cleared"
			m.statusTime = time.Now()
			return m, clearStatusAfterDelay()
		}
		// Extract unique tags from entries and todos
		m.availableTags = helpers.ExtractUniqueTagsFromAll(m.entries, m.todos)
		if len(m.availableTags) == 0 {
			m.statusMsg = "No tags found"
			m.statusTime = time.Now()
			return m, clearStatusAfterDelay()
		}
		// Open tag filter view
		m.tagFilterInput.Reset()
		m.tagFilterInput.Focus()
		m.autocompleteTag = ""
		m.view = "tag_filter"
		return m, textarea.Blink
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
		filtered := helpers.FilterEntriesByTags(m.entries, m.filterTags)

		if m.selectedEntry >= 0 && m.selectedEntry < len(filtered) {
			// Need to get the sorted entry (newest first)
			sorted := helpers.SortEntriesForDisplay(filtered)
			m.viewingEntry = sorted[m.selectedEntry]
			m.scrollOffset = 0 // Reset scroll when opening entry
			m.view = "view_entry"
			// Load todos so we can display them in the entry view
			return m, m.loadTodos()
		}
		return m, nil
	}
	return m, nil
}
