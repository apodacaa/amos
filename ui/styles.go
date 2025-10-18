package ui

import "github.com/charmbracelet/lipgloss"

// Brutalist color palette
var (
	subtleColor = lipgloss.AdaptiveColor{Light: "#000000", Dark: "#FFFFFF"}
	accentColor = lipgloss.AdaptiveColor{Light: "#000000", Dark: "#00FF00"}
	mutedColor  = lipgloss.AdaptiveColor{Light: "#666666", Dark: "#888888"}
)

// GetContainerStyle returns a container style sized to terminal dimensions
func GetContainerStyle(width, height int) lipgloss.Style {
	return lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(subtleColor).
		Padding(1, 2).
		Width(width - 2).  // Account for border
		Height(height - 2) // Account for border
}

// GetTitleStyle returns a title style sized to container width
func GetTitleStyle(width int) lipgloss.Style {
	return lipgloss.NewStyle().
		Bold(true).
		Foreground(accentColor).
		Align(lipgloss.Center).
		Width(width - 8) // Account for container padding + border
}

// Text styles
var (
	menuItemStyle = lipgloss.NewStyle().
			Foreground(subtleColor).
			MarginTop(0).
			MarginBottom(0)

	keyStyle = lipgloss.NewStyle().
			Foreground(accentColor).
			Bold(true)

	helpStyle = lipgloss.NewStyle().
			Foreground(mutedColor).
			Italic(true).
			MarginTop(1)
)
