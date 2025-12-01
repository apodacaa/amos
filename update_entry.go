package main

import (
	tea "github.com/charmbracelet/bubbletea"
)

// handleEntryKeys processes keyboard input (entry form view)
func (m Model) handleEntryKeys(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "ctrl+c":
		// Only Ctrl+C quits from entry form, not 'q' (user might type words with 'q')
		return m, tea.Quit
	case "esc":
		// Check if showing confirmation
		if m.confirmingExit {
			// User pressed Esc again - discard changes and exit
			m.textarea.Blur()
			m.confirmingExit = false
			m.statusMsg = ""
			m.hasUnsaved = false

			// If editing, return to view_entry; if creating, return to entries list
			if m.editingMode {
				m.editingMode = false
				m.originalEntryID = ""
				m.viewingEntry = m.currentEntry
				m.view = "view_entry"
				return m, nil
			}

			m.view = "entries"
			return m, m.loadEntriesAndTodos() // Reload to show any changes
		}

		// Check for unsaved changes
		currentContent := m.textarea.Value()
		if m.hasUnsaved && currentContent != m.savedContent {
			// Show confirmation prompt
			m.confirmingExit = true
			m.statusMsg = "Unsaved changes - Press Esc again to discard or Ctrl+S to save"
			return m, nil
		}

		// No unsaved changes, safe to exit
		m.textarea.Blur()
		m.confirmingExit = false

		// If editing, return to view_entry; if creating, return to entries list
		if m.editingMode {
			m.editingMode = false
			m.originalEntryID = ""
			m.viewingEntry = m.currentEntry
			m.view = "view_entry"
			return m, nil
		}

		m.view = "entries"
		return m, m.loadEntriesAndTodos() // Reload to show any changes

	case "ctrl+s":
		// Save entry
		m.confirmingExit = false // Clear confirmation if showing
		return m, m.saveEntry()

	default:
		// If confirming exit and user starts typing, cancel confirmation
		if m.confirmingExit {
			m.confirmingExit = false
			m.statusMsg = ""
		}

		// Let textarea handle the key
		var cmd tea.Cmd
		m.textarea, cmd = m.textarea.Update(msg)

		// Mark as having unsaved changes
		m.hasUnsaved = true

		return m, cmd
	}
}
