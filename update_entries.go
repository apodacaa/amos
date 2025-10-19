package main

import (
	"time"

	"github.com/apodacaa/amos/internal/helpers"
	"github.com/apodacaa/amos/internal/models"
	"github.com/charmbracelet/bubbles/textarea"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/google/uuid"
)

// handleEntriesListKeys processes keyboard input (entries list view)
func (m Model) handleEntriesListKeys(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	// If confirming delete, only allow 'd' to proceed or anything else to cancel
	if m.confirmingDelete {
		switch msg.String() {
		case "d":
			// User pressed 'd' again - proceed with delete
			return m, m.deleteEntry()
		default:
			// Any other key cancels delete confirmation
			m.confirmingDelete = false
			m.statusMsg = ""
			return m, nil
		}
	}

	// Normal navigation (not confirming delete)
	switch msg.String() {
	case "q", "ctrl+c":
		return m, tea.Quit
	case "esc":
		m.view = "dashboard"
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
	case "t":
		// Jump to todo list
		m.view = "todos"
		m.selectedTodo = 0
		// Todos already loaded (we load them when entering entries view)
		return m, nil
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
	case "@":
		// Open tag picker (or clear filter if already filtering)
		if m.filterTag != "" {
			// Clear filter
			m.filterTag = ""
			m.statusMsg = "✓ Filter cleared"
			return m, nil
		}
		// Extract unique tags from all entries
		m.availableTags = helpers.ExtractUniqueTags(m.entries)
		if len(m.availableTags) == 0 {
			m.statusMsg = "No tags found in entries"
			return m, nil
		}
		m.selectedTag = 0
		m.view = "tag_picker"
		return m, nil
	case "j", "down":
		if m.selectedEntry < len(m.entries)-1 {
			m.selectedEntry++
		}
		return m, nil
	case "k", "up":
		if m.selectedEntry > 0 {
			m.selectedEntry--
		}
		return m, nil
	case "d":
		// First 'd' press - show confirmation
		m.confirmingDelete = true
		m.statusMsg = "⚠ Delete entry? Press 'd' again to confirm, or any other key to cancel"
		return m, nil
	case "enter":
		// Open selected entry for read-only viewing
		// Apply filter if active (same logic as UI)
		filtered := helpers.FilterEntriesByTag(m.entries, m.filterTag)

		if m.selectedEntry >= 0 && m.selectedEntry < len(filtered) {
			// Need to get the sorted entry (newest first)
			sorted := helpers.SortEntriesForDisplay(filtered)
			m.viewingEntry = sorted[m.selectedEntry]
			m.view = "view_entry"
			// Load todos so we can display them in the entry view
			return m, m.loadTodos()
		}
		return m, nil
	}
	return m, nil
}
