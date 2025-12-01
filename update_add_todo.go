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
		// Check if showing confirmation
		if m.confirmingExit {
			// User pressed Esc again - discard changes and exit
			m.todoInput.Blur()
			m.confirmingExit = false
			m.statusMsg = ""
			m.hasUnsaved = false

			// If editing, return to view_todo; if creating, return to todos list
			if m.editingMode {
				m.editingMode = false
				m.originalTodoID = ""
				m.viewingTodo = m.currentTodo
				m.view = "view_todo"
				return m, nil
			}

			m.editingMode = false
			m.originalTodoID = ""
			m.view = "todos"
			return m, m.loadEntriesAndTodos() // Reload to show any changes
		}

		// Check for unsaved changes
		currentContent := strings.TrimSpace(m.todoInput.Value())
		savedContent := strings.TrimSpace(m.savedContent)
		if m.hasUnsaved && currentContent != savedContent {
			// Show confirmation prompt
			m.confirmingExit = true
			m.statusMsg = "Unsaved changes - Press Esc again to discard or Enter to save"
			return m, nil
		}

		// No content, safe to exit
		m.todoInput.Blur()
		m.confirmingExit = false

		// If editing, return to view_todo; if creating, return to todos list
		if m.editingMode {
			m.editingMode = false
			m.originalTodoID = ""
			m.viewingTodo = m.currentTodo
			m.view = "view_todo"
			return m, nil
		}

		m.editingMode = false
		m.originalTodoID = ""
		m.view = "todos"
		return m, m.loadEntriesAndTodos() // Reload to show any changes

	case "enter":
		title := strings.TrimSpace(m.todoInput.Value())
		if title == "" {
			m.statusMsg = "Todo title cannot be empty"
			return m, nil
		}

		// Set title and extract tags
		m.currentTodo.Title = title
		m.currentTodo.Tags = helpers.ExtractTags(title)

		// If editing, preserve existing fields; if creating, set as standalone
		if !m.editingMode {
			m.currentTodo.EntryID = nil // Standalone todo (no entry link)
		}
		// Note: When editing, currentTodo already has ID, Status, CreatedAt, EntryID, Position

		// Save current todo
		cmd := m.saveTodo()

		// If editing, clear edit mode and return to todo view
		// If creating, reset for next todo (power mode - rapid entry)
		if m.editingMode {
			m.editingMode = false
			m.originalTodoID = ""
			m.todoInput.Blur()
			m.viewingTodo = m.currentTodo // Update viewing todo with saved changes
			m.view = "view_todo"
		} else {
			// Reset for next todo (rapid entry mode)
			m.currentTodo = models.Todo{
				ID:        m.generateID(),
				Status:    "open",
				CreatedAt: time.Now(),
			}
			m.todoInput.Reset()
		}
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
