package main

import (
	"github.com/apodacaa/amos/internal/models"
	tea "github.com/charmbracelet/bubbletea"
)

// handleTodosListKeys processes keyboard input (todos list view)
func (m Model) handleTodosListKeys(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "q", "ctrl+c":
		return m, tea.Quit
	case "esc":
		// No need to save - already saved immediately
		m.view = "dashboard"
		m.statusMsg = ""
		return m, nil
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
			// Sort same way as display: open before done, then by position, then newest first
			sorted := make([]models.Todo, len(m.todos))
			copy(sorted, m.todos)
			for i := 0; i < len(sorted)-1; i++ {
				for j := i + 1; j < len(sorted); j++ {
					if sorted[i].Status == "done" && sorted[j].Status == "open" {
						sorted[i], sorted[j] = sorted[j], sorted[i]
					} else if sorted[i].Status == sorted[j].Status {
						if sorted[i].Position != sorted[j].Position {
							if sorted[j].Position < sorted[i].Position {
								sorted[i], sorted[j] = sorted[j], sorted[i]
							}
						} else {
							if sorted[j].CreatedAt.After(sorted[i].CreatedAt) {
								sorted[i], sorted[j] = sorted[j], sorted[i]
							}
						}
					}
				}
			}

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
