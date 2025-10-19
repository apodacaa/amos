package ui

import (
	"github.com/charmbracelet/bubbles/textarea"
	"github.com/charmbracelet/lipgloss"
)

// RenderEntryForm renders the entry editing form
func RenderEntryForm(width, height int, ta textarea.Model, statusMsg string) string {
	containerStyle := GetFullScreenBox(width, height)
	titleStyle := GetTitleStyle(width)

	title := titleStyle.Render("NEW ENTRY")

	help := FormatHelpLeft(width, "ctrl+s", "save", "esc", "exit")

	// Add status message if present
	status := ""
	if statusMsg != "" {
		statusStyle := lipgloss.NewStyle().Foreground(mutedColor)
		status = "\n" + statusStyle.Render(statusMsg)
	}

	content := title + "\n\n" +
		ta.View() + "\n\n" +
		help + status

	return containerStyle.Render(content)
}
