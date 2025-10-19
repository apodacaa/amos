package ui

import (
	"github.com/charmbracelet/bubbles/textarea"
	"github.com/charmbracelet/lipgloss"
)

// RenderAddTodoForm renders the standalone todo creation form
func RenderAddTodoForm(width, height int, todoInput textarea.Model, statusMsg string) string {
	container := GetContainerStyle(width, height)
	title := GetTitleStyle(width).Render("Add Todo")

	// Help text
	helpStyle := lipgloss.NewStyle().
		Foreground(mutedColor).
		Italic(true)
	help := helpStyle.Render("enter: save • esc: cancel")

	// Status message (if present)
	status := ""
	if statusMsg != "" {
		statusStyle := lipgloss.NewStyle().
			Foreground(mutedColor).
			Italic(true)
		status = "\n" + statusStyle.Render(statusMsg)
	}

	// Combine sections
	content := lipgloss.JoinVertical(
		lipgloss.Left,
		title,
		"",
		todoInput.View(),
		status,
		"",
		help,
	)

	return container.Render(content)
}
