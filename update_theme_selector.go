package main

import (
	tea "github.com/charmbracelet/bubbletea"

	"github.com/apodacaa/amos/ui"
)

// handleThemeSelectorKeys processes keyboard input (theme selector modal)
func (m Model) handleThemeSelectorKeys(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "q", "ctrl+c":
		return m, tea.Quit
	case "esc":
		// Cancel and return to previous view
		m.view = m.previousView
		m.statusMsg = ""
		return m, nil
	case "up", "k":
		// Move selection up
		if m.selectedTheme > 0 {
			m.selectedTheme--
		}
		return m, nil
	case "down", "j":
		// Move selection down
		themes := ui.ListThemes()
		if m.selectedTheme < len(themes)-1 {
			m.selectedTheme++
		}
		return m, nil
	case "enter":
		// Apply selected theme
		themes := ui.ListThemes()
		if m.selectedTheme >= 0 && m.selectedTheme < len(themes) {
			selectedTheme := themes[m.selectedTheme]
			m.currentTheme = selectedTheme

			// Update system config
			m.sysConfig.Theme = selectedTheme.Name

			// Save config and return to previous view
			m.view = m.previousView
			return m, saveConfigCmd(m.sysConfig)
		}
		return m, nil
	}

	return m, nil
}
