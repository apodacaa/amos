package main

import (
	tea "github.com/charmbracelet/bubbletea"
)

// handleViewEntryKeys processes keyboard input (view entry - read-only)
func (m Model) handleViewEntryKeys(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "q", "ctrl+c":
		return m, tea.Quit
	case "esc":
		// Go back to entries list
		m.view = "entries"
		return m, nil
	case "t":
		// Jump to todo list (todos already loaded from entry view)
		m.view = "todos"
		m.selectedTodo = 0
		// If somehow todos aren't loaded, load them
		if len(m.todos) == 0 {
			return m, m.loadTodos()
		}
		return m, nil
	}
	return m, nil
}
