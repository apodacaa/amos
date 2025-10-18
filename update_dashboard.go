package main

import (
	"time"

	"github.com/apodacaa/amos/internal/models"
	"github.com/charmbracelet/bubbles/textarea"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/google/uuid"
)

// handleKeyPress processes keyboard input (dashboard view)
func (m Model) handleKeyPress(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "q", "ctrl+c":
		return m, tea.Quit
	case "n":
		// Create new entry
		m.view = "entry"
		m.currentEntry = models.Entry{
			ID:        uuid.New().String(),
			Timestamp: time.Now(),
		}
		m.textarea.Reset()
		m.textarea.Focus()
		m.statusMsg = ""
		m.hasUnsaved = false
		m.savedContent = ""
		m.confirmingExit = false
		return m, textarea.Blink
	case "e":
		// View entries list
		m.view = "entries"
		m.selectedEntry = 0
		return m, m.loadEntries()
	case "t":
		// View todos list
		m.view = "todos"
		m.selectedTodo = 0
		return m, m.loadTodos()
	case "esc":
		m.view = "dashboard"
	}
	return m, nil
}
