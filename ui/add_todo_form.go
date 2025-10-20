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
		status = statusStyle.Render(statusMsg)
	}

	// Build main content (everything except help)
	mainContent := lipgloss.JoinVertical(
		lipgloss.Left,
		title,
		"",
		ti.View(),
	)
	if status != "" {
		mainContent = lipgloss.JoinVertical(lipgloss.Left, mainContent, "", status)
	}

	// Calculate how much vertical space to add to push help to bottom
	mainLines := lipgloss.Height(mainContent)
	helpLines := 1
	availableSpace := height - 4
	padding := availableSpace - mainLines - helpLines
	if padding < 0 {
		padding = 0
	}

	// Add padding and help
	content := mainContent
	if padding > 0 {
		content += "\n" + lipgloss.NewStyle().Height(padding).Render("")
	}
	content += "\n" + help

	return box.Render(content)
}
