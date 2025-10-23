package main

import (
	"strings"
	"time"

	"github.com/apodacaa/amos/internal/helpers"
	"github.com/apodacaa/amos/internal/models"
	tea "github.com/charmbracelet/bubbletea"
)

// handleAddTodoKeys processes keyboard input (add standalone todo form)
func (m Model) handleAddTodoKeys(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "ctrl+c":
		return m, tea.Quit
	case "esc":
		// Exit to dashboard (discard any unsaved input)
		m.view = "dashboard"
		m.todoInput.Blur()
		m.confirmingExit = false
		m.statusMsg = ""
		m.hasUnsaved = false
		return m, nil

	case "enter":
		// Save todo and start a new one (power mode - rapid entry)
		title := strings.TrimSpace(m.todoInput.Value())
		if title == "" {
			m.statusMsg = "⚠ Todo title cannot be empty"
			return m, nil
		}

		// Set title and extract tags
		m.currentTodo.Title = title
		m.currentTodo.Tags = helpers.ExtractTags(title)
		m.currentTodo.EntryID = nil // Standalone todo (no entry link)

		// Save current todo
		cmd := m.saveTodo()

		// Reset for next todo (rapid entry mode)
		m.currentTodo = models.Todo{
			ID:        m.generateID(),
			Status:    "open",
			CreatedAt: time.Now(),
		}
		m.todoInput.Reset()
		m.hasUnsaved = false

		return m, cmd
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
