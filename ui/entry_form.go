package ui

import (
	"github.com/charmbracelet/bubbles/textarea"
)

// RenderEntryForm renders the entry editing form
func RenderEntryForm(width, height int, ta textarea.Model, statusMsg string, hasUnsaved bool) string {
	return RenderFormView(width, height, ta.View(), statusMsg, hasUnsaved, "ctrl+s", "save", "esc", "cancel")
}
