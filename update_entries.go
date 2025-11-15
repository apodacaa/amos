package main

import (
	"strconv"
	"time"

	"github.com/apodacaa/amos/internal/helpers"
	"github.com/apodacaa/amos/ui"
	tea "github.com/charmbracelet/bubbletea"
)

// handleEntriesListKeys processes keyboard input (entries list view)
func (m Model) handleEntriesListKeys(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "q", "ctrl+c":
		return m, tea.Quit
	case "n":
		// Check if cancelling deletion first
		if m.deleteConfirmPending {
			m.deleteConfirmPending = false
			m.statusMsg = "Cancelled"
			m.statusTime = time.Now()
			return m, clearStatusAfterDelay()
		}
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
	case "s":
		// Open theme selector
		m.previousView = m.view
		m.view = "theme_selector"
		// Set selected theme to current theme
		themes := ui.ListThemes()
		for i, theme := range themes {
			if theme.Name == m.currentTheme.Name {
				m.selectedTheme = i
				break
			}
		}
		return m, nil
	case "?":
		// Open help page
		m.previousView = m.view
		m.scrollOffset = 0 // Reset scroll when opening help
		m.view = "help"
		return m, nil
	case "/":
		// Open unified filter input with current filter pre-populated
		return m.openUnifiedFilter("entries")
	case "c":
		// Clear all filters
		return m.clearFilters()
	case "j", "down":
		// Apply filters to get displayed list
		filtered := helpers.FilterEntriesByDateRange(m.entries, m.filterDate)
		filtered = helpers.FilterEntriesByTags(filtered, m.filterTags)

		if m.selectedEntry < len(filtered)-1 {
			m.selectedEntry++
		}
		return m, nil
	case "k", "up":
		if m.selectedEntry > 0 {
			m.selectedEntry--
		}
		return m, nil
	case "d":
		// Toggle mark for deletion
		filtered := helpers.FilterEntriesByDateRange(m.entries, m.filterDate)
		filtered = helpers.FilterEntriesByTags(filtered, m.filterTags)

		if m.selectedEntry >= 0 && m.selectedEntry < len(filtered) {
			sorted := helpers.SortEntriesForDisplay(filtered)
			selectedEntry := sorted[m.selectedEntry]

			// Toggle mark
			if _, marked := m.markedForDeletion[selectedEntry.ID]; marked {
				// Unmark
				delete(m.markedForDeletion, selectedEntry.ID)
				m.statusMsg = "Unmarked"
			} else {
				// Mark
				m.markedForDeletion[selectedEntry.ID] = "entry"
				// Count total marked items
				entryCount := 0
				todoCount := 0
				for _, itemType := range m.markedForDeletion {
					if itemType == "entry" {
						entryCount++
					} else {
						todoCount++
					}
				}
				if entryCount > 0 && todoCount > 0 {
					m.statusMsg = strconv.Itoa(entryCount) + " entries, " + strconv.Itoa(todoCount) + " todos marked. Press $ to delete."
				} else if entryCount > 0 {
					m.statusMsg = strconv.Itoa(entryCount) + " entries marked. Press $ to delete."
				} else {
					m.statusMsg = strconv.Itoa(todoCount) + " todos marked. Press $ to delete."
				}
			}
			m.statusTime = time.Now()
			return m, clearStatusAfterDelay()
		}
		return m, nil
	case "y":
		// Confirm deletion (if pending)
		if m.deleteConfirmPending {
			m.deleteConfirmPending = false
			m.statusMsg = ""
			return m, m.deleteMarkedCmd()
		}

	case "$":
		// Execute deletion of marked items
		if len(m.markedForDeletion) == 0 {
			m.statusMsg = "No items marked for deletion"
			m.statusTime = time.Now()
			return m, clearStatusAfterDelay()
		}

		// Show confirmation with y/n prompt
		// Count entries, todos, and affected linked todos
		entryCount := 0
		todoCount := 0
		linkedTodoCount := 0

		for id, itemType := range m.markedForDeletion {
			if itemType == "entry" {
				entryCount++
				// Find the entry and count its linked todos
				for _, e := range m.entries {
					if e.ID == id {
						linkedTodoCount += len(e.TodoIDs)
						break
					}
				}
			} else {
				todoCount++
			}
		}

		// Build confirmation message
		var msg string
		if entryCount > 0 && todoCount > 0 {
			msg = "Delete " + strconv.Itoa(entryCount) + " entries"
			if linkedTodoCount > 0 {
				msg += " (" + strconv.Itoa(linkedTodoCount) + " linked todos)"
			}
			msg += " and " + strconv.Itoa(todoCount) + " todos? ([yes]/no)"
		} else if entryCount > 0 {
			msg = "Delete " + strconv.Itoa(entryCount) + " entries"
			if linkedTodoCount > 0 {
				msg += " (" + strconv.Itoa(linkedTodoCount) + " linked todos)"
			}
			msg += "? ([yes]/no)"
		} else {
			msg = "Delete " + strconv.Itoa(todoCount) + " todos? ([yes]/no)"
		}

		m.deleteConfirmPending = true
		m.statusMsg = msg
		m.statusTime = time.Now()
		return m, nil

	case "enter":
		// Check if confirming deletion first
		if m.deleteConfirmPending {
			m.deleteConfirmPending = false
			m.statusMsg = ""
			return m, m.deleteMarkedCmd()
		}

		// Open selected entry for read-only viewing
		// Apply filters (same logic as UI)
		filtered := helpers.FilterEntriesByDateRange(m.entries, m.filterDate)
		filtered = helpers.FilterEntriesByTags(filtered, m.filterTags)

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
