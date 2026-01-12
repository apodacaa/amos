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
			// User pressed Esc again - DISCARD unsaved changes, finalize last saved version
			m.textarea.Blur()
			m.confirmingExit = false
			m.statusMsg = ""
			m.hasUnsaved = false

			// FINALIZE: Extract todos from LAST SAVED VERSION (m.currentEntry.Body)
			// Unsaved textarea content is discarded (standard TUI behavior)
			return m, m.finalizeAndExit()
		}

		// Check for unsaved changes
		currentContent := m.textarea.Value()
		if m.hasUnsaved && currentContent != m.savedContent {
			// Show confirmation prompt
			m.confirmingExit = true
			m.statusMsg = "Unsaved changes - Press Esc again to discard or Ctrl+S to save"
			return m, nil
		}

		// No unsaved changes - finalize current content and exit
		m.textarea.Blur()
		m.confirmingExit = false

		// FINALIZE: Extract todos from current saved state
		return m, m.finalizeAndExit()

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
