package main

import (
	"github.com/apodacaa/amos/internal/storage"
	tea "github.com/charmbracelet/bubbletea"
)

// handleWorkspaceSelectorKeys processes keyboard input (workspace selector modal)
func (m Model) handleWorkspaceSelectorKeys(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
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
		if m.selectedWorkspace > 0 {
			m.selectedWorkspace--
		}
		return m, nil
	case "down", "j":
		// Move selection down
		if m.selectedWorkspace < len(m.sysConfig.Workspaces)-1 {
			m.selectedWorkspace++
		}
		return m, nil
	case "enter":
		// Apply selected workspace
		workspaces := m.sysConfig.Workspaces
		if m.selectedWorkspace >= 0 && m.selectedWorkspace < len(workspaces) {
			selected := workspaces[m.selectedWorkspace]

			// Skip if already on this workspace
			if selected.Name == m.sysConfig.ActiveWorkspace {
				m.view = m.previousView
				return m, nil
			}

			// Update config
			m.sysConfig.ActiveWorkspace = selected.Name

			// Return to previous view and switch workspace
			m.view = m.previousView
			return m, tea.Batch(
				switchWorkspaceCmd(storage.Workspace{Name: selected.Name, DataDir: selected.DataDir}),
				saveConfigCmd(m.sysConfig),
			)
		}
		return m, nil
	}

	return m, nil
}
