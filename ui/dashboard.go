package ui

import (
	"strings"

	"github.com/apodacaa/amos/internal/models"
	"github.com/charmbracelet/lipgloss"
)

// RenderDashboard renders the main dashboard view
func RenderDashboard(width, height int, entries []models.Entry, todos []models.Todo) string {
	// Massive ASCII art title - centered
	title := lipgloss.NewStyle().
		Foreground(accentColor).
		Bold(true).
		Width(width - 8).
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

	// Help text - centered with bold keys
	help := FormatHelp(width,
		"n", "new entry",
		"a", "add todo",
		"t", "todos",
		"e", "entries",
		"q", "quit",
	)

	// Content - vertically centered in box
	boxContent := lipgloss.JoinVertical(
		lipgloss.Center,
		title,
		"",
		help,
	)

	// Full screen box that fills most of terminal
	boxStyle := GetFullScreenBox(width, height).
		AlignVertical(lipgloss.Center)

	return boxStyle.Render(boxContent)
}
