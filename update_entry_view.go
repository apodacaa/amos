package main

import (
	tea "github.com/charmbracelet/bubbletea"
)

// handleViewEntryKeys processes keyboard input (view entry - read-only)
func (m Model) handleViewEntryKeys(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
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
	case "a":
		// Add standalone todo (using shared helper)
		return m.handleAddTodo()
	case "e":
		// Jump to entries list (explicit navigation)
		m.view = "entries"
		m.selectedEntry = 0
		return m, m.loadEntriesAndTodos()
	case "t":
		// Jump to todo list (explicit navigation, load both todos and entries)
		m.view = "todos"
		m.selectedTodo = 0
		return m, m.loadEntriesAndTodos()
	}
	return m, nil
}
