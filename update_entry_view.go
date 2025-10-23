package main

import (
	"github.com/apodacaa/amos/internal/helpers"
	tea "github.com/charmbracelet/bubbletea"
)

// handleViewEntryKeys processes keyboard input (view entry - read-only)
func (m Model) handleViewEntryKeys(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
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
	case "a":
		// Add standalone todo (using shared helper)
		return m.handleAddTodo()
	case "e":
		// Jump to entries list (explicit navigation)
		m.view = "entries"
		m.selectedEntry = 0
		m.statusMsg = "" // Clear status message when changing views
		return m, m.loadEntriesAndTodos()
	case "t":
		// Jump to todo list (explicit navigation, load both todos and entries)
		m.view = "todos"
		m.selectedTodo = 0
		m.statusMsg = "" // Clear status message when changing views
		return m, m.loadEntriesAndTodos()
	case "j", "down":
		// Navigate to next entry (newer to older, same as entry list)
		// Apply filter and sort (same as entry list view)
		filtered := helpers.FilterEntriesByTags(m.entries, m.filterTags)
		sorted := helpers.SortEntriesForDisplay(filtered)

		if len(sorted) > 0 {
			// Find current entry index in sorted list
			currentIdx := -1
			for i, entry := range sorted {
				if entry.ID == m.viewingEntry.ID {
					currentIdx = i
					break
				}
			}

			// Move down (to next entry, which is older)
			if currentIdx >= 0 && currentIdx < len(sorted)-1 {
				m.selectedEntry = currentIdx + 1
				m.viewingEntry = sorted[m.selectedEntry]
				m.scrollOffset = 0 // Reset scroll when switching entries
			}
		}
		return m, nil
	case "k", "up":
		// Navigate to previous entry (older to newer, same as entry list)
		// Apply filter and sort (same as entry list view)
		filtered := helpers.FilterEntriesByTags(m.entries, m.filterTags)
		sorted := helpers.SortEntriesForDisplay(filtered)

		if len(sorted) > 0 {
			// Find current entry index in sorted list
			currentIdx := -1
			for i, entry := range sorted {
				if entry.ID == m.viewingEntry.ID {
					currentIdx = i
					break
				}
			}

			// Move up (to previous entry, which is newer)
			if currentIdx > 0 {
				m.selectedEntry = currentIdx - 1
				m.viewingEntry = sorted[m.selectedEntry]
				m.scrollOffset = 0 // Reset scroll when switching entries
			}
		}
		return m, nil
	case "u":
		// Scroll down in current entry (u is above j, j goes down)
		m.scrollOffset++
		return m, nil
	case "i":
		// Scroll up in current entry (i is above k, k goes up)
		if m.scrollOffset > 0 {
			m.scrollOffset--
		}
		return m, nil
	}
	return m, nil
}
