package ui

import (
	"fmt"

	"github.com/charmbracelet/lipgloss"
)

// RenderTagPicker renders the tag selection view
func RenderTagPicker(width, height int, tags []string, selectedIdx int) string {
	container := GetFullScreenBox(width, height)
	title := GetTitleStyle(width).Render("Filter by Tag")

	// Build tag list
	var listItems []string

	if len(tags) == 0 {
		emptyStyle := lipgloss.NewStyle().
			Foreground(mutedColor).
			Width(width - 4).
			Align(lipgloss.Center)
		listItems = append(listItems, emptyStyle.Render("No tags found"))
	} else {
		for i, tag := range tags {
			var line string
			if i == selectedIdx {
				// Selected item (highlighted)
				line = fmt.Sprintf("▶ %s", tag)
				line = lipgloss.NewStyle().
					Foreground(accentColor).
					Bold(true).
					Render(line)
			} else {
				// Unselected item
				line = fmt.Sprintf("  %s", tag)
				line = lipgloss.NewStyle().
					Foreground(subtleColor).
					Render(line)
			}
			listItems = append(listItems, line)
		}
	}

	// Join list items
	list := lipgloss.JoinVertical(lipgloss.Left, listItems...)

	// Help text at bottom
	help := FormatHelpLeft(width,
		"enter", "select",
		"esc", "cancel",
	)

	// Build main content (everything except help)
	mainContent := lipgloss.JoinVertical(
		lipgloss.Left,
		title,
		"",
		list,
	)

	// Calculate how much vertical space to add to push help to bottom
	mainLines := lipgloss.Height(mainContent)
	helpLines := 1
	availableSpace := height - 4
	padding := availableSpace - mainLines - helpLines
	if padding < 0 {
		padding = 0
	}

	// Add padding and help
	content := mainContent
	if padding > 0 {
		content += "\n" + lipgloss.NewStyle().Height(padding).Render("")
	}
	content += "\n" + help

	return container.Render(content)
}
