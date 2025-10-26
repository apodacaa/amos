package ui

import (
	"strings"

	"github.com/charmbracelet/bubbles/textarea"
	"github.com/charmbracelet/lipgloss"
)

// RenderAddTodoForm renders the standalone todo creation form
func RenderAddTodoForm(width, height int, ti textarea.Model, statusMsg string, hasUnsaved bool) string {
	// Header
	header := RenderHeader(width, "enter", "save", "esc", "cancel")

	// Footer with unsaved indicator in statusMsg
	footerTitle := "Add Todo"
	footerStatus := statusMsg
	if footerStatus == "" && hasUnsaved {
		footerStatus = "Unsaved"
	}
	footer := RenderFooter(width, footerTitle, footerStatus)

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
