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
		status = statusStyle.Render(statusMsg)
	}

	// Build main content - help anchored to bottom
	mainContent := lipgloss.JoinVertical(
		lipgloss.Left,
		title,
		"",
		ta.View(),
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

	return containerStyle.Render(fullContent)
}
