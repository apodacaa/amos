package main

import (
	tea "github.com/charmbracelet/bubbletea"

	"github.com/apodacaa/amos/internal/helpers"
)

// handleDateFilterKeys handles key presses in the date filter view
func (m Model) handleDateFilterKeys(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	presets := helpers.GetDatePresets()

	switch msg.String() {
	case "q", "ctrl+c":
		return m, tea.Quit

	case "esc":
		// Cancel and go to dashboard
		m.view = "dashboard"
		return m, nil

	case "n":
		// Create new entry (using shared helper)
		return m.handleNewEntry()

	case "a":
		// Add standalone todo (using shared helper)
		return m.handleAddTodo()

	case "e":
		// Jump to entry list
		m.view = "entries"
		m.selectedEntry = 0
		return m, m.loadEntriesAndTodos()

	case "t":
		// Jump to todo list
		m.view = "todos"
		m.selectedTodo = 0
		return m, m.loadEntriesAndTodos()

	case "j", "down":
		// Move down
		if m.selectedDatePreset < len(presets)-1 {
			m.selectedDatePreset++
		}
		return m, nil

	case "k", "up":
		// Move up
		if m.selectedDatePreset > 0 {
			m.selectedDatePreset--
		}
		return m, nil

	case "enter":
		// Apply selected preset
		selectedPreset := presets[m.selectedDatePreset]
		m.filterDate = selectedPreset

		// Return to filter view
		m.view = "filter_view"

		return m, nil
	}

	return m, nil
}
