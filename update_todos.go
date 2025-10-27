package main

import (
	"github.com/apodacaa/amos/internal/helpers"
	tea "github.com/charmbracelet/bubbletea"
)

// handleTodosListKeys processes keyboard input (todos list view)
func (m Model) handleTodosListKeys(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "q", "ctrl+c":
		return m, tea.Quit
	case "n":
		// Create new entry (using shared helper)
		return m.handleNewEntry()
	case "e":
		// Jump to entry list (explicit navigation)
		m.view = "entries"
		m.selectedEntry = 0
		m.statusMsg = "" // Clear status message when changing views
		return m, m.loadEntriesAndTodos()
	case "a":
		// Add standalone todo (using shared helper)
		return m.handleAddTodo()
	case "/":
		// Open unified filter input with current filter pre-populated
		return m.openUnifiedFilter("todos")
	case "c":
		// Clear all filters
		return m.clearFilters()
	case "j", "down":
		// Apply filters to get the displayed list (same as UI)
		filtered := helpers.FilterTodosByDateRange(m.displayTodos, m.filterDate)
		filtered = helpers.FilterTodosByTags(filtered, m.filterTags)

		if m.selectedTodo < len(filtered)-1 {
			m.selectedTodo++
		}
		return m, nil
	case "k", "up":
		if m.selectedTodo > 0 {
			m.selectedTodo--
		}
		return m, nil
	case "r":
		// Refresh - reload todos to re-sort
		return m, m.loadTodos()
	case " ":
		// Cycle todo status: open → next → done → open (save immediately, no re-sort)
		// Use filtered displayTodos to keep selection stable
		filtered := helpers.FilterTodosByDateRange(m.displayTodos, m.filterDate)
		filtered = helpers.FilterTodosByTags(filtered, m.filterTags)
		if m.selectedTodo >= 0 && m.selectedTodo < len(filtered) {
			// Get the todo from filtered list (current display order)
			todo := filtered[m.selectedTodo]

			// Cycle status: open → next → done → open
			switch todo.Status {
			case "open":
				todo.Status = "next"
				m.statusMsg = "→ Next"
			case "next":
				todo.Status = "done"
				m.statusMsg = "✓ Done"
			case "done":
				todo.Status = "open"
				m.statusMsg = "○ Open"
			default:
				// Unknown status, set to open
				todo.Status = "open"
				m.statusMsg = "○ Open"
			}

			// Update in m.todos array and displayTodos (find by ID)
			// We can't update displayTodos[m.selectedTodo] directly because
			// we're working with a filtered view
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
	return m, nil
}
