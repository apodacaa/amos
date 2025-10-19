package main

import (
	"time"

	"github.com/apodacaa/amos/internal/models"
	"github.com/charmbracelet/bubbles/textarea"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/google/uuid"
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
	case "n":
		// Create new entry
		m.view = "entry"
		m.currentEntry = models.Entry{
			ID:        uuid.New().String(),
			Timestamp: time.Now(),
		}
		m.textarea.Reset()
		m.textarea.Focus()
		m.hasUnsaved = false
		m.savedContent = ""
		m.statusMsg = ""
		return m, textarea.Blink
	case "a":
		// Add standalone todo
		m.view = "add_todo"
		m.currentTodo = models.Todo{
			ID:        uuid.New().String(),
			Status:    "open",
			Position:  0,
			CreatedAt: time.Now(),
		}
		m.todoInput.Reset()
		m.todoInput.Focus()
		m.statusMsg = ""
		return m, textarea.Blink
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
