package ui

import (
	"github.com/charmbracelet/bubbles/textarea"
	"github.com/charmbracelet/lipgloss"
)

// RenderEntryForm renders the entry editing form
func RenderEntryForm(width, height int, ta textarea.Model, statusMsg string) string {
	containerStyle := GetFullScreenBox(width, height)
	titleStyle := GetTitleStyle(width)

	help := FormatHelpLeft(width, "ctrl+s", "save", "esc", "cancel")

	title := titleStyle.Render("NEW ENTRY")

	mainContent := lipgloss.JoinVertical(
		lipgloss.Left,
		title,
		"",
		ta.View(),
	)

	mainLines := lipgloss.Height(mainContent)
	helpLines := 1
	statusLines := 0

	// Add status message if present
	var statusRendered string
	if statusMsg != "" {
		statusStyle := lipgloss.NewStyle().Foreground(subtleColor)
		statusRendered = statusStyle.Render(statusMsg)
		statusLines = 1
	}

	availableSpace := height - 4
	padding := availableSpace - mainLines - helpLines - statusLines
	if padding < 0 {
		padding = 0
	}

	content := mainContent
	if padding > 0 {
		content += "\n" + lipgloss.NewStyle().Height(padding).Render("")
	}
	if statusMsg != "" {
		content += "\n" + statusRendered
	}
	content += "\n" + help

	return containerStyle.Render(content)
}
