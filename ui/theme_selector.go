package ui

import (
	"strings"

	"github.com/charmbracelet/lipgloss"
)

// RenderThemeSelector renders a centered modal for theme selection
func RenderThemeSelector(width, height int, theme Theme, selectedIdx int) string {
	themes := ListThemes()

	// Modal dimensions
	modalWidth := 60
	if modalWidth > width-4 {
		modalWidth = width - 4
	}

	// Build theme list
	var themeItems []string
	for i, t := range themes {
		prefix := "  "
		if i == selectedIdx {
			prefix = "> "
		}

		// Theme name in bold
		name := lipgloss.NewStyle().Bold(true).Render(t.Name)

		// Theme description
		desc := "  " + t.Description

		item := prefix + name + "\n" + desc

		// Highlight selected item
		if i == selectedIdx {
			itemStyle := lipgloss.NewStyle().
				Width(modalWidth-4).
				Padding(0, 1)

			// Use theme's selection style if available
			if theme.Colors.Selection != "" {
				itemStyle = itemStyle.
					Background(theme.Colors.Selection).
					Foreground(theme.Colors.Foreground)
			} else {
				itemStyle = itemStyle.Reverse(true)
			}

			item = itemStyle.Render(item)
		} else {
			item = lipgloss.NewStyle().
				Width(modalWidth-4).
				Padding(0, 1).
				Render(item)
		}

		themeItems = append(themeItems, item)
	}

	themeList := strings.Join(themeItems, "\n\n")

	// Title
	title := lipgloss.NewStyle().
		Bold(true).
		Width(modalWidth - 4).
		Align(lipgloss.Center).
		Render("Select Theme")

	// Help text
	help := FormatHelpLeft(modalWidth-4,
		"j/k", "navigate",
		"enter", "select",
		"esc", "cancel",
	)

	// Build modal content
	content := lipgloss.JoinVertical(
		lipgloss.Left,
		title,
		"",
		themeList,
		"",
		help,
	)

	// Modal box
	modalBox := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(theme.Colors.Primary).
		Padding(1, 2).
		Width(modalWidth).
		Render(content)

	// Center modal on screen
	modalHeight := lipgloss.Height(modalBox)
	verticalPadding := (height - modalHeight) / 2
	if verticalPadding < 0 {
		verticalPadding = 0
	}

	horizontalPadding := (width - modalWidth) / 2
	if horizontalPadding < 0 {
		horizontalPadding = 0
	}

	// Create centered modal
	centered := lipgloss.Place(
		width,
		height,
		lipgloss.Center,
		lipgloss.Center,
		modalBox,
	)

	return centered
}
