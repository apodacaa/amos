package ui

import (
	"strings"

	"github.com/charmbracelet/bubbles/textarea"
	"github.com/charmbracelet/lipgloss"
)

// RenderEntryForm renders the entry editing form
func RenderEntryForm(width, height int, ta textarea.Model, statusMsg string) string {
	// Header
	header := RenderHeader(width, "ctrl+s", "save", "esc", "cancel")

	// Footer
	footer := RenderFooter(width, "New Entry", statusMsg)

	// Calculate padding for content area
	contentHeight := height - 2 // header + footer
	textareaLines := lipgloss.Height(ta.View())
	padding := contentHeight - textareaLines
	if padding < 0 {
		padding = 0
	}

	// Build full view
	content := header + "\n" + ta.View()
	if padding > 0 {
		content += strings.Repeat("\n", padding)
	}
	content += "\n" + footer

	return content
}
