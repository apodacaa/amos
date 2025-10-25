package main

import (
	"github.com/charmbracelet/bubbles/textarea"
	tea "github.com/charmbracelet/bubbletea"

	"github.com/apodacaa/amos/internal/helpers"
)

// handleFilterViewKeys handles key presses in the unified filter view
func (m Model) handleFilterViewKeys(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
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

	case "@":
		// Open tag filter input
		m.availableTags = helpers.ExtractUniqueTagsFromAll(m.entries, m.todos)
		if len(m.availableTags) == 0 {
			// No tags available, stay in filter view
			return m, nil
		}
		m.tagFilterInput.Reset()
		m.tagFilterInput.Focus()
		m.autocompleteTag = ""
		m.view = "tag_filter"
		return m, textarea.Blink

	case "d":
		// Open date filter menu
		m.selectedDatePreset = 0
		m.view = "date_filter"
		return m, nil

	case "c":
		// Clear all filters
		m.filterTags = []string{}
		m.filterDate = ""
		return m, nil

	case "enter":
		// Apply filters and return to appropriate list
		m.view = m.filterContext
		// Reset selection based on context
		if m.filterContext == "entries" {
			m.selectedEntry = 0
		} else if m.filterContext == "todos" {
			m.selectedTodo = 0
		}
		return m, nil
	}

	return m, nil
}
