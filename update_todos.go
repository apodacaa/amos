package main

import (
	"strconv"
	"time"

	"github.com/apodacaa/amos/internal/helpers"
	"github.com/apodacaa/amos/ui"
	tea "github.com/charmbracelet/bubbletea"
)

// handleTodosListKeys processes keyboard input (todos list view)
func (m Model) handleTodosListKeys(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	// Check for update dismissal key 'u' (only if update notice is shown)
	if m.updateAvailable && !m.updateDismissed && msg.String() == "u" {
		m.updateDismissed = true
		return m, nil
	}

	switch msg.String() {
	case "q", "ctrl+c":
		return m, tea.Quit
	case "o":
		// Set todo to open
		return m.setTodoStatus("open", "Todo set to open")
	case "p":
		// Set todo to next (p for priority)
		return m.setTodoStatus("next", "Todo set to next")
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
	case "x":
		// Set todo to done
		return m.setTodoStatus("done", "Todo marked done")
	case "e":
		// Jump to entry list (explicit navigation)
		m.view = "entries"
		m.selectedEntry = 0
		m.statusMsg = "" // Clear status message when changing views
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
		return m.openUnifiedFilter("todos")
	case "c":
		// Clear all filters
		return m.clearFilters()
	case "z":
		// Toggle visibility of completed todos
		m.sysConfig.HideCompletedTodos = !m.sysConfig.HideCompletedTodos
		if m.sysConfig.HideCompletedTodos {
			m.statusMsg = "Completed todos hidden"
		} else {
			m.statusMsg = "Showing all todos"
		}
		m.statusTime = time.Now()
		m.selectedTodo = 0 // Reset selection to avoid out-of-bounds
		return m, tea.Batch(saveConfigCmd(m.sysConfig), clearStatusAfterDelay())
	case "j", "down":
		// Apply filters to get the displayed list (same as UI)
		filtered := helpers.ApplyTodoFilters(m.displayTodos, m.filterDate, m.filterTags, m.sysConfig.HideCompletedTodos)

		if m.selectedTodo < len(filtered)-1 {
			m.selectedTodo++
		}
		return m, nil
	case "k", "up":
		if m.selectedTodo > 0 {
			m.selectedTodo--
		}
		return m, nil
	case "d":
		// Toggle mark for deletion
		filtered := helpers.ApplyTodoFilters(m.displayTodos, m.filterDate, m.filterTags, m.sysConfig.HideCompletedTodos)

		if m.selectedTodo >= 0 && m.selectedTodo < len(filtered) {
			selectedTodo := filtered[m.selectedTodo]

			// Toggle mark
			if _, marked := m.markedForDeletion[selectedTodo.ID]; marked {
				// Unmark
				delete(m.markedForDeletion, selectedTodo.ID)
				m.statusMsg = "Unmarked"
			} else {
				// Mark (allow marking even entry-linked todos - they won't be deleted)
				m.markedForDeletion[selectedTodo.ID] = "todo"
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

	case "enter":
		// Check if confirming deletion first
		if m.deleteConfirmPending {
			m.deleteConfirmPending = false
			m.statusMsg = ""
			return m, m.deleteMarkedCmd()
		}

		// Open selected todo for read-only viewing
		// Apply filters (same logic as UI)
		filtered := helpers.ApplyTodoFilters(m.displayTodos, m.filterDate, m.filterTags, m.sysConfig.HideCompletedTodos)

		if m.selectedTodo >= 0 && m.selectedTodo < len(filtered) {
			m.viewingTodo = filtered[m.selectedTodo]
			m.scrollOffset = 0 // Reset scroll when opening todo
			m.view = "view_todo"
			// Entries already loaded for linked entry display
			return m, nil
		}
		return m, nil
	}
	return m, nil
}

// setTodoStatus sets the status of the selected todo and saves immediately
func (m Model) setTodoStatus(status string, statusMsg string) (tea.Model, tea.Cmd) {
	// Apply filters to get the displayed list
	filtered := helpers.ApplyTodoFilters(m.displayTodos, m.filterDate, m.filterTags, m.sysConfig.HideCompletedTodos)

	if m.selectedTodo >= 0 && m.selectedTodo < len(filtered) {
		todo := filtered[m.selectedTodo]
		todo.Status = status
		m.statusMsg = statusMsg
		m.statusTime = time.Now()

		// Update in m.todos array and displayTodos (find by ID)
		for i := range m.todos {
			if m.todos[i].ID == todo.ID {
				m.todos[i].Status = todo.Status
				break
			}
		}
		for i := range m.displayTodos {
			if m.displayTodos[i].ID == todo.ID {
				m.displayTodos[i].Status = todo.Status
				break
			}
		}

		// Save immediately and start timer to clear status
		return m, tea.Batch(m.toggleTodoImmediate(todo), clearStatusAfterDelay())
	}
	return m, nil
}
