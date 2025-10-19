package main

import (
	tea "github.com/charmbracelet/bubbletea"
)

// handleKeyPress processes keyboard input (dashboard view)
func (m Model) handleKeyPress(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "q", "ctrl+c":
		return m, tea.Quit
	case "n":
		// Create new entry (using shared helper)
		return m.handleNewEntry()
	case "e":
		// View entries list (load both entries and todos for stats)
		m.view = "entries"
		m.selectedEntry = 0
		return m, m.loadEntriesAndTodos()
	case "t":
		// View todos list
		m.view = "todos"
		m.selectedTodo = 0
		return m, m.loadTodos()
	case "a":
		// Add standalone todo (using shared helper)
		return m.handleAddTodo()
	case "esc":
		m.view = "dashboard"
	}
	return m, nil
}
