package ui

import "github.com/charmbracelet/lipgloss"

// Brutalist color palette
// Gainsboro (#DBDCE8), Metallic Silver (#AAA3B4), Purple Taupe (#463848),
// Deep Tuscan Red (#6E3F52), Mauve Taupe (#976775)
var (
	// Primary text and borders - Purple Taupe (dark) / Gainsboro (light)
	subtleColor = lipgloss.AdaptiveColor{Light: "#463848", Dark: "#DBDCE8"}

	// Accent/highlight - Deep Tuscan Red (both themes)
	accentColor = lipgloss.AdaptiveColor{Light: "#6E3F52", Dark: "#976775"}

	// Muted/help text - Metallic Silver (light) / Mauve Taupe (dark)
	mutedColor = lipgloss.AdaptiveColor{Light: "#AAA3B4", Dark: "#976775"}
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

// Textarea style helpers
func GetTextareaStyle() lipgloss.Style {
	return lipgloss.NewStyle().Foreground(subtleColor)
}

func GetPlaceholderStyle() lipgloss.Style {
	return lipgloss.NewStyle().Foreground(mutedColor)
}

func GetPromptStyle() lipgloss.Style {
	return lipgloss.NewStyle().Foreground(accentColor)
}

func GetTextStyle() lipgloss.Style {
	return lipgloss.NewStyle().Foreground(subtleColor)
}
