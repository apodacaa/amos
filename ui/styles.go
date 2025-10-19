package ui

import "github.com/charmbracelet/lipgloss"

// Brutalist color palette - Pure monochrome concrete
var (
	// Primary text and borders - Dark gray (light) / Light gray (dark)
	subtleColor = lipgloss.AdaptiveColor{Light: "#404040", Dark: "#CCCCCC"}

	// Accent/highlight - Pure black/white for maximum contrast
	accentColor = lipgloss.AdaptiveColor{Light: "#000000", Dark: "#FFFFFF"}

	// Muted/help text - Mid gray
	mutedColor = lipgloss.AdaptiveColor{Light: "#808080", Dark: "#666666"}
)

// GetFullScreenBox returns a box that fills most of the terminal with consistent styling
func GetFullScreenBox(width, height int) lipgloss.Style {
	return lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(accentColor).
		Padding(1, 2).
		Width(width - 2).  // Minimal margin for border
		Height(height - 2) // Minimal margin for border
}

// GetTitleStyle returns a title style sized to container width
func GetTitleStyle(width int) lipgloss.Style {
	return lipgloss.NewStyle().
		Bold(true).
		Foreground(accentColor).
		Width(width - 8).
		Align(lipgloss.Center)
}

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

// FormatHelp formats help text with bold keys (reverse colors for impact)
// Centered alignment for dashboard monument aesthetics
// Example: FormatHelp(width, "n", "new entry", "a", "add todo")
func FormatHelp(width int, keyDescPairs ...string) string {
	return formatHelpWithAlign(width, lipgloss.Center, keyDescPairs...)
}

// FormatHelpLeft formats help text with bold keys (reverse colors for impact)
// Left-aligned for utility views (honest functional UI)
// Example: FormatHelpLeft(width, "n", "new entry", "a", "add todo")
func FormatHelpLeft(width int, keyDescPairs ...string) string {
	return formatHelpWithAlign(width, lipgloss.Left, keyDescPairs...)
}

// formatHelpWithAlign is the shared implementation for help text formatting
func formatHelpWithAlign(width int, align lipgloss.Position, keyDescPairs ...string) string {
	var parts []string

	// Maximum contrast: invert accent colors (white on black / black on white)
	keyFg := lipgloss.AdaptiveColor{Light: "#FFFFFF", Dark: "#000000"}
	keyBg := lipgloss.AdaptiveColor{Light: "#000000", Dark: "#FFFFFF"}

	keyStyle := lipgloss.NewStyle().
		Foreground(keyFg).
		Background(keyBg).
		Bold(true).
		Padding(0, 1). // Add small padding for readability
		Inline(true)   // Keep inline to prevent breaking

	descStyle := lipgloss.NewStyle().
		Foreground(subtleColor).
		Inline(true) // Keep inline to prevent breaking

	for i := 0; i < len(keyDescPairs); i += 2 {
		if i+1 < len(keyDescPairs) {
			key := keyStyle.Render(keyDescPairs[i])
			desc := descStyle.Render(" " + keyDescPairs[i+1])
			// Combine key and desc as single unit with inline wrapper
			pair := lipgloss.NewStyle().Inline(true).Render(key + desc)
			parts = append(parts, pair)
		}
	}

	// Join parts with spacing
	result := ""
	for i, part := range parts {
		result += part
		if i < len(parts)-1 {
			result += "  " // Two spaces between items
		}
	}

	return lipgloss.NewStyle().
		Foreground(mutedColor).
		Width(width - 8).
		Align(align).
		Inline(true). // Prevent wrapping
		Render(result)
}
