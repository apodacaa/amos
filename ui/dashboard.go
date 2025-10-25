package ui

import (
	"fmt"
	"strings"

	"github.com/apodacaa/amos/internal/helpers"
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

	// Weekly stats (90 days = ~13 weeks)
	weekStats := helpers.AggregateByWeek(entries, todos, 13)

	// Calculate available height for graph (total - header - footer - title)
	titleLines := 8 // ASCII art lines
	availableHeight := height - 2 - titleLines
	if availableHeight < 15 {
		availableHeight = 15 // Minimum for readable graph
	}

	statsSection := RenderLineGraph(weekStats, width, availableHeight)

	// Calculate content area (height - header - footer)
	contentHeight := height - 2 // 1 for header, 1 for footer
	titleHeight := lipgloss.Height(title)
	statsLines := lipgloss.Height(statsSection)
	mainContentLines := titleHeight + statsLines + 1 // +1 for spacing

	padding := contentHeight - mainContentLines
	if padding < 0 {
		padding = 0
	}

	// Build main content: title, spacing, stats
	mainContent := title + "\n" + statsSection

	// Build full view
	content := header + "\n" + mainContent
	if padding > 0 {
		content += strings.Repeat("\n", padding)
	}
	content += "\n" + footer

	return content
}
