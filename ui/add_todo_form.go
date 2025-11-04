package ui

import (
	"github.com/charmbracelet/bubbles/textarea"
)

// RenderAddTodoForm renders the standalone todo creation form
func RenderAddTodoForm(width, height int, ti textarea.Model, statusMsg string, hasUnsaved bool) string {
	return RenderFormView(width, height, ti.View(), statusMsg, hasUnsaved, "enter", "save", "esc", "cancel")
}
