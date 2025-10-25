package ui

import (
	"strings"

	"github.com/apodacaa/amos/internal/helpers"
	"github.com/charmbracelet/lipgloss"
)

// RenderDateFilter renders the date preset selection menu
func RenderDateFilter(width, height int, selectedIdx int) string {
	container := GetFullScreenBox(width, height)
	presets := helpers.GetDatePresets()

	// Title
	titleText := "Date Filter"
	title := GetTitleStyle(width).Render(titleText)

	// Menu items
	menuItems := []string{}
	for i, preset := range presets {
		label := helpers.FormatDatePreset(preset)
		if label == "" {
			label = "all"
		}

		if i == selectedIdx {
			// Selected item (subtle color)
			item := lipgloss.NewStyle().
				Foreground(subtleColor).
				Render("> " + label)
			menuItems = append(menuItems, item)
		} else {
			// Unselected item (muted)
			item := lipgloss.NewStyle().
				Foreground(mutedColor).
				Render("  " + label)
			menuItems = append(menuItems, item)
		}
	}

	menu := strings.Join(menuItems, "\n")

	// Help text - full navigation available
	help := FormatHelpLeft(width,
		"n", "new entry",
		"a", "add todo",
		"j/k", "navigate",
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
		menu,
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
