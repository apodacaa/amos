package main

import (
	"time"

	tea "github.com/charmbracelet/bubbletea"
)

// handleTodosListKeys processes keyboard input (todos list view)
func (m Model) handleTodosListKeys(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
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
	case "e":
		// Jump to entry list (explicit navigation)
		m.view = "entries"
		m.selectedEntry = 0
		m.statusMsg = "" // Clear status message when changing views
		return m, m.loadEntriesAndTodos()
	case "a":
		// Add standalone todo (using shared helper)
		return m.handleAddTodo()
	case "j", "down":
		if m.selectedTodo < len(m.displayTodos)-1 {
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
		// Use displayTodos to keep selection stable
		if m.selectedTodo >= 0 && m.selectedTodo < len(m.displayTodos) {
			// Get the todo from displayTodos (current display order)
			todo := m.displayTodos[m.selectedTodo]

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
			m.statusTime = time.Now()

			// Update in displayTodos (in place, no re-sort)
			m.displayTodos[m.selectedTodo].Status = todo.Status

			// Update in m.todos array (find by ID)
			for i := range m.todos {
				if m.todos[i].ID == todo.ID {
					m.todos[i].Status = todo.Status
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
