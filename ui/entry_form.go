package ui

import (
	"github.com/charmbracelet/bubbles/textarea"
	"github.com/charmbracelet/lipgloss"
)

// RenderEntryForm renders the entry editing form
func RenderEntryForm(width, height int, ta textarea.Model) string {
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
	availableSpace := height - 4
	padding := availableSpace - mainLines - helpLines
	if padding < 0 {
		padding = 0
	}

	content := mainContent
	if padding > 0 {
		content += "\n" + lipgloss.NewStyle().Height(padding).Render("")
	}
	content += "\n" + help

	return containerStyle.Render(content)
}
