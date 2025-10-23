package ui

import (
	"github.com/charmbracelet/bubbles/textarea"
	"github.com/charmbracelet/lipgloss"
)

// RenderAddTodoForm renders the standalone todo creation form
func RenderAddTodoForm(width, height int, ti textarea.Model) string {
	box := GetFullScreenBox(width, height)
	titleStyle := GetTitleStyle(width)

	help := FormatHelpLeft(width, "enter", "save", "esc", "cancel")

	title := titleStyle.Render("ADD TODO")

	mainContent := lipgloss.JoinVertical(
		lipgloss.Left,
		title,
		"",
		ti.View(),
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

	return box.Render(content)
}
