package main

import (
	"strconv"
	"time"

	"github.com/apodacaa/amos/internal/helpers"
	"github.com/charmbracelet/bubbles/textarea"
	tea "github.com/charmbracelet/bubbletea"
)

// handleViewTodoKeys processes keyboard input (view todo - read-only)
func (m Model) handleViewTodoKeys(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "q", "ctrl+c":
		return m, tea.Quit
	case "enter":
		// Check if confirming deletion first
		if m.deleteConfirmPending {
			m.deleteConfirmPending = false
			m.statusMsg = ""
			m.view = "todos" // Return to todo list after deletion
			return m, m.deleteMarkedCmd()
		}

		// Go back to todo list
		m.view = "todos"
		m.statusMsg = "" // Clear status message when changing views
		return m, nil

	case "esc":
		// Go back to todo list
		m.view = "todos"
		m.statusMsg = "" // Clear status message when changing views
		return m, nil
	case "i":
		// Edit current todo
		return m.handleEditTodo()
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
	case " ":
		// Check if confirming/cancelling deletion first
		if m.deleteConfirmPending {
			// Space doesn't do anything during confirmation (use y/n/enter)
			return m, nil
		}

		// Cycle todo status: open → next → done → open (save immediately)
		// Update the viewing todo
		switch m.viewingTodo.Status {
		case "open":
			m.viewingTodo.Status = "next"
			m.statusMsg = "Todo marked as next"
		case "next":
			m.viewingTodo.Status = "done"
			m.statusMsg = "Todo marked as done"
		case "done":
			m.viewingTodo.Status = "open"
			m.statusMsg = "Todo marked as open"
		default:
			// Unknown status, set to open
			m.viewingTodo.Status = "open"
			m.statusMsg = "Todo marked as open"
		}
		m.statusTime = time.Now()

		// Update in m.todos array and displayTodos (find by ID)
		for i := range m.todos {
			if m.todos[i].ID == m.viewingTodo.ID {
				m.todos[i].Status = m.viewingTodo.Status
				break
			}
		}
		for i := range m.displayTodos {
			if m.displayTodos[i].ID == m.viewingTodo.ID {
				m.displayTodos[i].Status = m.viewingTodo.Status
				break
			}
		}

		// Save immediately and start timer to clear status
		return m, tea.Batch(m.toggleTodoImmediate(m.viewingTodo), clearStatusAfterDelay())
	case "j", "down":
		// Navigate to next todo (open → next → done order)
		// Apply filters and sort (same as todo list view)
		filtered := helpers.FilterTodosByDateRange(m.displayTodos, m.filterDate)
		filtered = helpers.FilterTodosByTags(filtered, m.filterTags)

		if len(filtered) > 0 {
			// Find current todo index in filtered list
			currentIdx := -1
			for i, todo := range filtered {
				if todo.ID == m.viewingTodo.ID {
					currentIdx = i
					break
				}
			}

			// Move down (to next todo)
			if currentIdx >= 0 && currentIdx < len(filtered)-1 {
				m.selectedTodo = currentIdx + 1
				m.viewingTodo = filtered[m.selectedTodo]
				m.scrollOffset = 0 // Reset scroll when switching todos
			}
		}
		return m, nil
	case "k", "up":
		// Navigate to previous todo (open → next → done order)
		// Apply filters and sort (same as todo list view)
		filtered := helpers.FilterTodosByDateRange(m.displayTodos, m.filterDate)
		filtered = helpers.FilterTodosByTags(filtered, m.filterTags)

		if len(filtered) > 0 {
			// Find current todo index in filtered list
			currentIdx := -1
			for i, todo := range filtered {
				if todo.ID == m.viewingTodo.ID {
					currentIdx = i
					break
				}
			}

			// Move up (to previous todo)
			if currentIdx > 0 {
				m.selectedTodo = currentIdx - 1
				m.viewingTodo = filtered[m.selectedTodo]
				m.scrollOffset = 0 // Reset scroll when switching todos
			}
		}
		return m, nil
	case "b":
		// Scroll backward (up) one page
		availableHeight := m.height - 3
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
		availableHeight := m.height - 3
		if availableHeight < 5 {
			availableHeight = 5
		}
		m.scrollOffset += availableHeight
		// UI layer will clamp to maxOffset, so we don't need to calculate it here
		return m, nil
	case "d":
		// Mark/unmark current todo for deletion
		if _, marked := m.markedForDeletion[m.viewingTodo.ID]; marked {
			delete(m.markedForDeletion, m.viewingTodo.ID)
			m.statusMsg = "Unmarked"
		} else {
			m.markedForDeletion[m.viewingTodo.ID] = "todo"

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

			// Build status message
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

	case "y":
		// Confirm deletion (if pending)
		if m.deleteConfirmPending {
			m.deleteConfirmPending = false
			m.statusMsg = ""
			m.view = "todos" // Return to todo list after deletion
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
	}
	return m, nil
}

// handleEditTodo transitions from todo view to editing it
func (m Model) handleEditTodo() (tea.Model, tea.Cmd) {
	// Copy viewing todo to current todo (preserves all fields)
	m.currentTodo = m.viewingTodo

	// Set editing mode flags
	m.editingMode = true
	m.originalTodoID = m.viewingTodo.ID

	// Load todo title into input
	m.todoInput.SetValue(m.viewingTodo.Title)
	m.todoInput.Focus()

	// Reset editing state - set savedContent AFTER SetValue to ensure exact match
	m.hasUnsaved = false
	m.savedContent = m.todoInput.Value() // Get actual value from textarea
	m.confirmingExit = false

	// Switch to add_todo form view (reused for editing)
	m.view = "add_todo"

	return m, textarea.Blink
}
