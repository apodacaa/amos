package main

import (
	"strings"

	"github.com/apodacaa/amos/internal/helpers"
	tea "github.com/charmbracelet/bubbletea"
)

// handleAddTodoKeys processes keyboard input (add standalone todo form)
func (m Model) handleAddTodoKeys(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "ctrl+c":
		return m, tea.Quit
	case "esc":
		// Cancel and go to dashboard
		m.view = "dashboard"
		m.statusMsg = ""
		return m, nil
	case "ctrl+s", "enter":
		// Save standalone todo
		title := strings.TrimSpace(m.todoInput.Value())
		if title == "" {
			m.statusMsg = "⚠ Todo title cannot be empty"
			return m, nil
		}

		// Set title and extract tags
		m.currentTodo.Title = title
		m.currentTodo.Tags = helpers.ExtractTags(title)
		m.currentTodo.EntryID = nil // Standalone todo (no entry link)

		// Save and return to dashboard or todo list
		return m, m.saveTodo()
	default:
		// Let all other keys pass through to textarea
		var cmd tea.Cmd
		m.todoInput, cmd = m.todoInput.Update(msg)
		return m, cmd
	}
}
