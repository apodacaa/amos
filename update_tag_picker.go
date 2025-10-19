package main

import (
	tea "github.com/charmbracelet/bubbletea"
)

// handleTagPickerKeys processes keyboard input (tag picker view)
func (m Model) handleTagPickerKeys(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "q", "ctrl+c":
		return m, tea.Quit
	case "esc":
		// Go back to entries without filtering
		m.view = "entries"
		m.statusMsg = "" // Clear any previous status message
		return m, nil
	case "j", "down":
		if m.selectedTag < len(m.availableTags)-1 {
			m.selectedTag++
		}
		return m, nil
	case "k", "up":
		if m.selectedTag > 0 {
			m.selectedTag--
		}
		return m, nil
	case "enter":
		// Apply filter and go back to entries
		if m.selectedTag >= 0 && m.selectedTag < len(m.availableTags) {
			m.filterTag = m.availableTags[m.selectedTag]
			m.view = "entries"
			m.selectedEntry = 0 // Reset selection
			m.statusMsg = ""    // Clear any previous status message
		}
		return m, nil
	}
	return m, nil
}
