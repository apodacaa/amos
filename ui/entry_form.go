package ui

import "github.com/charmbracelet/bubbles/textarea"

// RenderEntryForm renders the entry form view
func RenderEntryForm(width, height int, ta textarea.Model, statusMsg string) string {
	containerStyle := GetContainerStyle(width, height)
	titleStyle := GetTitleStyle(width)

	title := titleStyle.Render("NEW ENTRY")

	help := helpStyle.Render("Ctrl+S to save • Esc to exit")

	// Add status message if present
	status := ""
	if statusMsg != "" {
		status = "\n" + helpStyle.Render(statusMsg)
	}

	content := title + "\n\n" +
		ta.View() + "\n\n" +
		help + status

	return containerStyle.Render(content)
}
