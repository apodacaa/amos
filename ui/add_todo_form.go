package ui

import (
	"github.com/charmbracelet/bubbles/textarea"
	"github.com/charmbracelet/lipgloss"
)

// RenderAddTodoForm renders the standalone todo creation form
func RenderAddTodoForm(width, height int, ti textarea.Model, statusMsg string) string {
	box := GetFullScreenBox(width, height)
	titleStyle := GetTitleStyle(width)

	title := titleStyle.Render("ADD TODO")

	help := FormatHelpLeft(width, "enter", "save", "esc", "cancel")

	// Status message (if present)
	status := ""
	if statusMsg != "" {
		statusStyle := lipgloss.NewStyle().
			Foreground(mutedColor)
		status = "\n" + statusStyle.Render(statusMsg)
	}

	// Combine sections
	content := lipgloss.JoinVertical(
		lipgloss.Left,
		title,
		"",
		ti.View(),
		status,
		"",
		help,
	)

	return box.Render(content)
}
