package main

import (
	"strconv"
	"strings"
	"time"

	"github.com/apodacaa/amos/internal/helpers"
	"github.com/apodacaa/amos/ui"
	tea "github.com/charmbracelet/bubbletea"
)

// handleViewEntryKeys processes keyboard input (view entry - read-only)
func (m Model) handleViewEntryKeys(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	// Check for update dismissal key 'u' (only if update notice is shown)
	if m.updateAvailable && !m.updateDismissed && msg.String() == "u" {
		m.updateDismissed = true
		return m, nil
	}

	switch msg.String() {
	case "q", "ctrl+c":
		return m, tea.Quit
	case "esc":
		// Go back to entry list
		m.view = "entries"
		m.statusMsg = "" // Clear status message when changing views
		return m, nil
	case "i":
		// Edit current entry
		return m.handleEditEntry()
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
		// Apply filters and sort (same as entry list view)
		filtered := helpers.ApplyEntryFilters(m.entries, m.filterDate, m.filterTags)
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
		// Apply filters and sort (same as entry list view)
		filtered := helpers.ApplyEntryFilters(m.entries, m.filterDate, m.filterTags)
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
	case "b":
		// Scroll backward (up) one page
		// Calculate availableHeight matching ui/entry_view.go
		wrappedTitle := ui.WrapBodyText(m.viewingEntry.Title, m.width-8)
		titleLineCount := strings.Count(wrappedTitle, "\n") + 1
		baseReserved := 6 + titleLineCount
		if m.updateAvailable && !m.updateDismissed {
			baseReserved += 1
		}
		// Calculate todos height
		entryTodos := helpers.FilterTodosByEntry(m.todos, m.viewingEntry.ID)
		todosHeight := 0
		if len(entryTodos) > 0 {
			// 2 blank lines + title line + one line per todo
			todosHeight = 2 + 1 + len(entryTodos)
		}
		availableHeight := m.height - baseReserved - todosHeight
		if availableHeight < 5 {
			availableHeight = 5
		}
		m.scrollOffset -= availableHeight
		if m.scrollOffset < 0 {
			m.scrollOffset = 0
		}
		return m, nil
	case "f":
		// Scroll forward (down) one page
		// Calculate availableHeight matching ui/entry_view.go
		wrappedTitle := ui.WrapBodyText(m.viewingEntry.Title, m.width-8)
		titleLineCount := strings.Count(wrappedTitle, "\n") + 1
		baseReserved := 6 + titleLineCount
		if m.updateAvailable && !m.updateDismissed {
			baseReserved += 1
		}
		// Calculate todos height
		entryTodos := helpers.FilterTodosByEntry(m.todos, m.viewingEntry.ID)
		todosHeight := 0
		if len(entryTodos) > 0 {
			// 2 blank lines + title line + one line per todo
			todosHeight = 2 + 1 + len(entryTodos)
		}
		availableHeight := m.height - baseReserved - todosHeight
		if availableHeight < 5 {
			availableHeight = 5
		}

		// Calculate max offset using wrapped body
		wrappedBody := ui.WrapBodyText(m.viewingEntry.Body, m.width-8)
		bodyLines := strings.Split(wrappedBody, "\n")
		totalLines := len(bodyLines)

		maxOffset := totalLines - availableHeight
		if maxOffset < 0 {
			maxOffset = 0
		}

		// Increment and clamp to max
		m.scrollOffset += availableHeight
		if m.scrollOffset > maxOffset {
			m.scrollOffset = maxOffset
		}
		return m, nil
	case "d":
		// Toggle mark for deletion (for currently viewed entry)
		// Toggle mark
		if _, marked := m.markedForDeletion[m.viewingEntry.ID]; marked {
			// Unmark
			delete(m.markedForDeletion, m.viewingEntry.ID)
			m.statusMsg = "Unmarked"
		} else {
			// Mark
			m.markedForDeletion[m.viewingEntry.ID] = "entry"
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
	case "y", "enter":
		// Confirm deletion (if pending)
		if m.deleteConfirmPending {
			m.deleteConfirmPending = false
			m.statusMsg = ""
			m.view = "entries" // Return to entries list after deletion
			return m, m.deleteMarkedCmd()
		}

	case "$":
		// Execute deletion of marked items (same logic as entries list)
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
	}
	return m, nil
}

// handleEditEntry transitions from viewing an entry to editing it
func (m Model) handleEditEntry() (tea.Model, tea.Cmd) {
	// Copy viewing entry to current entry (preserves ID and all fields)
	m.currentEntry = m.viewingEntry

	// Set editing mode flags
	m.editingMode = true
	m.originalEntryID = m.viewingEntry.ID

	// Format content for editor: "title\n\nbody"
	content := m.viewingEntry.Title
	if m.viewingEntry.Body != "" {
		content += "\n\n" + m.viewingEntry.Body
	}

	return m, launchEditor(content)
}
