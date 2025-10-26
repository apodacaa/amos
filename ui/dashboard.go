package ui

import (
	"fmt"
	"strings"

	"github.com/apodacaa/amos/internal/models"
	"github.com/charmbracelet/lipgloss"
)

// RenderDashboard renders the main dashboard view
func RenderDashboard(width, height int, entries []models.Entry, todos []models.Todo) string {
	// Header
	header := RenderHeader(width, "n", "new", "a", "todo", "t", "todos", "e", "entries", "q", "quit")

	// Calculate stats for footer
	totalEntries := len(entries)
	openTodos := 0
	nextTodos := 0
	doneTodos := 0
	for _, todo := range todos {
		switch todo.Status {
		case "open":
			openTodos++
		case "next":
			nextTodos++
		case "done":
			doneTodos++
		}
	}

	// Footer with stats
	footerStats := fmt.Sprintf("%d entries, %d open todos, %d next todos, %d done todos", totalEntries, openTodos, nextTodos, doneTodos)
	footer := RenderFooter(width, "Dashboard", footerStats)

	// Massive ASCII art title - centered
	title := lipgloss.NewStyle().
		Foreground(accentColor).
		Bold(true).
		Width(width).
		Align(lipgloss.Center).
		Render(strings.Join([]string{
			"",
			"░█████╗░███╗░░░███╗░█████╗░░██████╗",
			"██╔══██╗████╗░████║██╔══██╗██╔════╝",
			"███████║██╔████╔██║██║░░██║╚█████╗░",
			"██╔══██║██║╚██╔╝██║██║░░██║░╚═══██╗",
			"██║░░██║██║░╚═╝░██║╚█████╔╝██████╔╝",
			"╚═╝░░╚═╝╚═╝░░░░░╚═╝░╚════╝░╚═════╝░",
			"",
		}, "\n"))

	// Calculate content area (height - header - footer)
	contentHeight := height - 2 // 1 for header, 1 for footer
	titleHeight := lipgloss.Height(title)

	// Calculate padding to center title vertically
	padding := (contentHeight - titleHeight) / 2
	if padding < 0 {
		padding = 0
	}

	// Build centered content
	var content strings.Builder
	content.WriteString(header)
	content.WriteString("\n")

	// Top padding
	if padding > 0 {
		content.WriteString(strings.Repeat("\n", padding))
	}

	// Centered title
	content.WriteString(title)

	// Bottom padding
	remainingPadding := contentHeight - titleHeight - padding
	if remainingPadding > 0 {
		content.WriteString(strings.Repeat("\n", remainingPadding))
	}

	content.WriteString("\n")
	content.WriteString(footer)

	return content.String()
}
