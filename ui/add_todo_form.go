package ui

import (
	"strings"

	"github.com/charmbracelet/bubbles/textarea"
	"github.com/charmbracelet/lipgloss"
)

// RenderAddTodoForm renders the standalone todo creation form
func RenderAddTodoForm(width, height int, ti textarea.Model, statusMsg string) string {
	// Header
	header := RenderHeader(width, "enter", "save", "esc", "cancel")

	// Footer
	footer := RenderFooter(width, "Add Todo", statusMsg)

	// Calculate padding for content area
	contentHeight := height - 2 // header + footer
	textareaLines := lipgloss.Height(ti.View())
	padding := contentHeight - textareaLines
	if padding < 0 {
		padding = 0
	}

	// Build full view
	content := header + "\n" + ti.View()
	if padding > 0 {
		content += strings.Repeat("\n", padding)
	}
	content += "\n" + footer

	return content
}
