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
		// Check if showing confirmation
		if m.confirmingExit {
			// User pressed Esc again - discard changes and exit to dashboard
			m.view = "dashboard"
			m.todoInput.Blur()
			m.confirmingExit = false
			m.statusMsg = ""
			m.hasUnsaved = false
			return m, nil
		}

		// Check for unsaved changes
		currentContent := strings.TrimSpace(m.todoInput.Value())
		if m.hasUnsaved && currentContent != "" {
			// Show confirmation prompt
			m.confirmingExit = true
			m.statusMsg = "⚠ Unsaved changes! Press Esc again to discard, or Ctrl+S to save"
			return m, nil
		}

		// No unsaved changes, safe to exit to dashboard
		m.view = "dashboard"
		m.todoInput.Blur()
		m.confirmingExit = false
		return m, nil

	case "ctrl+s":
		// Save standalone todo
		m.confirmingExit = false // Clear confirmation if showing
		title := strings.TrimSpace(m.todoInput.Value())
		if title == "" {
			m.statusMsg = "⚠ Todo title cannot be empty"
			return m, nil
		}

		// Set title and extract tags
		m.currentTodo.Title = title
		m.currentTodo.Tags = helpers.ExtractTags(title)
		m.currentTodo.EntryID = nil // Standalone todo (no entry link)

		// Save and return to dashboard
		return m, m.saveTodo()
	default:
		// If confirming exit and user starts typing, cancel confirmation
		if m.confirmingExit {
			m.confirmingExit = false
			m.statusMsg = ""
		}

		// Let all other keys pass through to textarea
		var cmd tea.Cmd
		m.todoInput, cmd = m.todoInput.Update(msg)

		// Mark as having unsaved changes
		m.hasUnsaved = true

		return m, cmd
	}
}
