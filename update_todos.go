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
	case "d":
		// Go to dashboard (explicit navigation)
		m.view = "dashboard"
		return m, nil
	case "n":
		// Create new entry (using shared helper)
		return m.handleNewEntry()
	case "e":
		// Jump to entry list (explicit navigation)
		m.view = "entries"
		m.selectedEntry = 0
		return m, m.loadEntriesAndTodos()
	case "a":
		// Add standalone todo (using shared helper)
		return m.handleAddTodo()
	case "j", "down":
		if m.selectedTodo < len(m.todos)-1 {
			m.selectedTodo++
		}
		return m, nil
	case "k", "up":
		if m.selectedTodo > 0 {
			m.selectedTodo--
		}
		return m, nil
	case "u":
		// Move todo down (lower priority)
		cmd := m.moveTodo("down")
		// Keep selection on the same todo (which will now be one position down)
		if m.selectedTodo < len(m.todos)-1 {
			m.selectedTodo++
		}
		return m, cmd
	case "i":
		// Move todo up (higher priority)
		cmd := m.moveTodo("up")
		// Keep selection on the same todo (which will now be one position up)
		if m.selectedTodo > 0 {
			m.selectedTodo--
		}
		return m, cmd
	case " ":
		// Toggle todo status (save immediately, no reload)
		// Need to sort todos same way as display to get the right one
		if m.selectedTodo >= 0 && m.selectedTodo < len(m.todos) {
			// Sort using helper (same logic as UI and commands)
			sorted := helpers.SortTodosForDisplay(m.todos)

			// Get the todo from the sorted list (matches display order)
			todo := sorted[m.selectedTodo]

			// Toggle status
			if todo.Status == "done" {
				todo.Status = "open"
				m.statusMsg = "✓ Reopened"
			} else {
				todo.Status = "done"
				m.statusMsg = "✓ Done"
			}

			// Update in m.todos array (find by ID)
			for i := range m.todos {
				if m.todos[i].ID == todo.ID {
					m.todos[i].Status = todo.Status
					break
				}
			}

			// Save immediately (async, don't wait)
			return m, m.toggleTodoImmediate(todo)
		}
		return m, nil
	case "s":
		// No longer needed, but keep for consistency (does nothing)
		return m, nil
	}
	return m, nil
}
