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

	help := FormatHelpLeft(width, "enter", "save", "ctrl+s", "save", "esc", "cancel")

	// Status message (if present)
	status := ""
	if statusMsg != "" {
		statusStyle := lipgloss.NewStyle().
			Foreground(mutedColor)
		status = statusStyle.Render(statusMsg)
	}

	// Build main content - help anchored to bottom
	mainContent := lipgloss.JoinVertical(
		lipgloss.Left,
		title,
		"",
		ti.View(),
	)
	if status != "" {
		mainContent = lipgloss.JoinVertical(lipgloss.Left, mainContent, "", status)
	}

	// Place content with help anchored to bottom
	fullContent := lipgloss.Place(
		width-4, height-4,
		lipgloss.Left, lipgloss.Top,
		mainContent,
	) + "\n" + help

	return box.Render(fullContent)
}
