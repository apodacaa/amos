package main

import (
	"github.com/apodacaa/amos/internal/helpers"
	tea "github.com/charmbracelet/bubbletea"
)

// handleViewTodoKeys processes keyboard input (view todo - read-only)
func (m Model) handleViewTodoKeys(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "q", "ctrl+c":
		return m, tea.Quit
	case "esc", "enter":
		// Go back to todo list
		m.view = "todos"
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
	case " ":
		// Cycle todo status: open → next → done → open (save immediately)
		// Update the viewing todo
		switch m.viewingTodo.Status {
		case "open":
			m.viewingTodo.Status = "next"
			m.statusMsg = "→ Next"
		case "next":
			m.viewingTodo.Status = "done"
			m.statusMsg = "✓ Done"
		case "done":
			m.viewingTodo.Status = "open"
			m.statusMsg = "○ Open"
		default:
			// Unknown status, set to open
			m.viewingTodo.Status = "open"
			m.statusMsg = "○ Open"
		}

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
	case "d":
		// Scroll down in current todo (d = down, standard TUI convention)
		m.scrollOffset++
		return m, nil
	case "u":
		// Scroll up in current todo (u = up, standard TUI convention)
		if m.scrollOffset > 0 {
			m.scrollOffset--
		}
		return m, nil
	}
	return m, nil
}
