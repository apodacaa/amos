package main

import (
	"github.com/apodacaa/amos/internal/models"
	tea "github.com/charmbracelet/bubbletea"
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
		if m.selectedEntry >= 0 && m.selectedEntry < len(m.entries) {
			// Need to get the sorted entry (newest first)
			sorted := make([]models.Entry, len(m.entries))
			copy(sorted, m.entries)
			// Sort by timestamp descending (newest first) - same as in RenderEntryList
			for i := 0; i < len(sorted)-1; i++ {
				for j := i + 1; j < len(sorted); j++ {
					if sorted[j].Timestamp.After(sorted[i].Timestamp) {
						sorted[i], sorted[j] = sorted[j], sorted[i]
					}
				}
			}
			m.viewingEntry = sorted[m.selectedEntry]
			m.view = "view_entry"
			// Load todos so we can display them in the entry view
			return m, m.loadTodos()
		}
		return m, nil
	}
	return m, nil
}
