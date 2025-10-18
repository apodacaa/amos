package ui

import (
	"fmt"
	"strings"

	"github.com/apodacaa/amos/internal/models"
	"github.com/charmbracelet/lipgloss"
)

// RenderEntryView renders a read-only view of an entry
func RenderEntryView(width, height int, entry models.Entry) string {
	container := GetContainerStyle(width, height)

	// Title
	titleStyle := lipgloss.NewStyle().
		Bold(true).
		Foreground(accentColor).
		Width(width - 8)
	title := titleStyle.Render(entry.Title)

	// Metadata line: date and tags
	metaStyle := lipgloss.NewStyle().
		Foreground(mutedColor).
		Italic(true)

	timestamp := entry.Timestamp.Format("2006-01-02 15:04")
	meta := timestamp

	if len(entry.Tags) > 0 {
		meta += " • " + strings.Join(entry.Tags, " ")
	}

	metadata := metaStyle.Render(meta)

	// Body
	bodyStyle := lipgloss.NewStyle().
		Foreground(subtleColor).
		Width(width - 8)
	body := bodyStyle.Render(entry.Body)

	// Todos section (if any)
	var todosSection string
	if len(entry.TodoIDs) > 0 {
		todosTitle := lipgloss.NewStyle().
			Bold(true).
			Foreground(accentColor).
			Render(fmt.Sprintf("Todos (%d)", len(entry.TodoIDs)))

		todosInfo := lipgloss.NewStyle().
			Foreground(mutedColor).
			Italic(true).
			Render("(Linked todos - view in todo list)")

		todosSection = "\n\n" + todosTitle + "\n" + todosInfo
	}

	// Help text
	helpStyle := lipgloss.NewStyle().
		Foreground(mutedColor).
		Italic(true)
	help := helpStyle.Render("esc: back to list • q: quit")

	// Combine all sections
	content := lipgloss.JoinVertical(
		lipgloss.Left,
		title,
		metadata,
		"",
		body,
		todosSection,
		"",
		help,
	)

	return container.Render(content)
}
