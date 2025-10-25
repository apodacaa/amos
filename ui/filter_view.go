package ui

import (
	"strings"

	"github.com/charmbracelet/lipgloss"
)

// RenderFilterView renders the unified filter management view
func RenderFilterView(width, height int, filterTags []string, filterDate string, dateLabel string) string {
	container := GetFullScreenBox(width, height)

	// Title
	titleText := "Filters"
	title := GetTitleStyle(width).Render(titleText)

	// Current filter status
	var filterStatus []string

	// Tags status
	tagsLine := "Tags: "
	if len(filterTags) > 0 {
		tagsLine += strings.Join(filterTags, " ")
	} else {
		tagsLine += "none"
	}
	filterStatus = append(filterStatus, tagsLine)

	// Date status
	dateLine := "Date: "
	if filterDate != "" && dateLabel != "" {
		dateLine += dateLabel
	} else {
		dateLine += "none"
	}
	filterStatus = append(filterStatus, dateLine)

	status := lipgloss.NewStyle().
		Foreground(subtleColor).
		Render(strings.Join(filterStatus, "\n"))

	// Help text - full navigation available
	help := FormatHelpLeft(width,
		"n", "new entry",
		"a", "add todo",
		"@", "tags",
		"d", "date",
		"c", "clear all",
		"enter", "apply",
		"e", "entries",
		"t", "todos",
		"esc", "cancel",
		"q", "quit",
	)

	// Build main content (everything except help)
	mainContent := lipgloss.JoinVertical(
		lipgloss.Left,
		title,
		"",
		status,
	)

	// Calculate how much vertical space to add to push help to bottom
	mainLines := strings.Count(mainContent, "\n") + 1
	helpLines := 1
	availableSpace := height - 4
	padding := availableSpace - mainLines - helpLines
	if padding < 0 {
		padding = 0
	}

	// Add padding and help
	content := mainContent
	if padding > 0 {
		content += strings.Repeat("\n", padding)
	}
	content += "\n" + help

	return container.Render(content)
}
