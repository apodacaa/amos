package main

import (
	tea "github.com/charmbracelet/bubbletea"
)

// handleTodosListKeys processes keyboard input (todos list view)
func (m Model) handleTodosListKeys(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "q", "ctrl+c":
		return m, tea.Quit
	case "esc":
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
		// Move todo up (higher priority)
		cmd := m.moveTodo("up")
		// Keep selection on the same todo (which will now be one position up)
		if m.selectedTodo > 0 {
			m.selectedTodo--
		}
		return m, cmd
	case "i":
		// Move todo down (lower priority)
		cmd := m.moveTodo("down")
		// Keep selection on the same todo (which will now be one position down)
		if m.selectedTodo < len(m.todos)-1 {
			m.selectedTodo++
		}
		return m, cmd
	case " ":
		// Toggle todo status
		return m, m.toggleTodo()
	}
	return m, nil
}
